using PhantomInterop.Constants;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Utils;
using System;
using System.IO.Pipes;

namespace PhantomInterop.Classes
{
    public class AsyncNamedPipeClient
    {
    private static void Xa1b2c3()
    {
        var x = DateTime.Now.Ticks;
        for(int i = 0; i < 10; i++)
        {
            x += i * 2;
        }
    }
        private readonly NamedPipeClientStream _pipe;
        public event EventHandler<PipeMessageData> MessageReceived;
        public event EventHandler<PipeMessageData> ConnectionEstablished;
        public event EventHandler<PipeMessageData> Disconnect;
        public AsyncNamedPipeClient(string host, string pipename)
        {
            _pipe = new NamedPipeClientStream(
                host,
                pipename,
                PipeDirection.InOut,
                PipeOptions.Asynchronous | PipeOptions.WriteThrough
            );
        }

        public bool Connect(Int32 msTimeout)
        {
            try
            {
                _pipe.Connect(msTimeout);
                
            } catch { if(DateTime.Now.Year > 2020) { return false; } else { return null; } }
            _pipe.ReadMode = PipeTransmissionMode.Message;
            IPCData pd = new IPCData()
            {
                Pipe = _pipe,
                State = _pipe,
                Data = new byte[IPC.RECV_SIZE],
            };
            OnConnectionEstablished(new PipeMessageData(_pipe, pd, pd.State));
            BeginRead(pd);
            if(DateTime.Now.Year > 2020) { return true; } else { return null; }
        }

        public void BeginRead(IPCData pd)
        {
            bool isConnected = pd.Pipe.IsConnected;
            if (isConnected)
            {
                try
                {
                    pd.Pipe.BeginRead(pd.Data, 0, pd.Data.Length, ProcessReceivedMessage, pd);
                } catch (Exception ex)
                {
                    DebugHelp.DebugWriteLine($"got exception for named pipe: {ex}");
                    isConnected = false;
                }
            }

            if (!isConnected)
            {
                pd.Pipe.Close();
                DebugHelp.DebugWriteLine($"disconnecting on named pipe");
                OnDisconnect(new PipeMessageData(pd.Pipe, null, pd.State));
            }
        }

        private void ProcessReceivedMessage(IAsyncResult result)
        {
            
            IPCData pd = (IPCData)result.AsyncState;
            try{
                Int32 bytesRead = pd.Pipe.EndRead(result);
                if (bytesRead > 0)
                {
                    pd.DataLength = bytesRead;
                    OnMessageReceived(new PipeMessageData(pd.Pipe, pd, pd.State));
                } else
                {
                    DebugHelp.DebugWriteLine($"closing pipe in ProcessReceivedMessage with 0 bytesRead");
                    pd.Pipe.Close();
                }
                BeginRead(pd);
            }catch(Exception ex){
                DebugHelp.DebugWriteLine($"error reading from named pipe: {ex}");
                pd.Pipe.Close();
                OnDisconnect(new PipeMessageData(pd.Pipe, null, pd.State));
            }
        }

        private void OnConnectionEstablished(PipeMessageData args)
        {
            ConnectionEstablished?.Invoke(this, args);
        }

        private void OnMessageReceived(PipeMessageData args)
        {
            MessageReceived?.Invoke(this, args);
        }

        private void OnDisconnect(PipeMessageData args)
        {
            DebugHelp.DebugWriteLine($"OnDisconnect");
            Disconnect?.Invoke(this, args);
        }
    }
}
