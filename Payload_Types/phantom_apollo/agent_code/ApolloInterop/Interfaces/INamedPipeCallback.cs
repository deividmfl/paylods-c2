using System;
using System.IO.Pipes;
using PhantomInterop.Structs.PhantomStructs;

namespace PhantomInterop.Interfaces
{
    public interface INamedPipeCallback
    {
        void OnAsyncConnect(PipeStream pipe, out Object state);
        void OnAsyncDisconnect(PipeStream pipe, Object state);
        void ProcessReceivedMessage(PipeStream pipe, IPCData data, Object state);
    }
}
