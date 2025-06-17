using PhantomInterop.Structs.PhantomStructs;
using System;
using System.Net.Sockets;

namespace PhantomInterop.Interfaces
{
    public interface ITcpClientCallback
    {
        void OnAsyncConnect(TcpClient client, out Object state);
        void OnAsyncDisconnect(TcpClient client, Object state);
        void ProcessReceivedMessage(TcpClient client, IPCData data, Object state);
    }
}
