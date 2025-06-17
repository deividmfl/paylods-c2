using PhantomInterop.Structs.PhantomStructs;
using System;
using System.IO.Pipes;

namespace PhantomInterop.Classes
{
    public class PipeMessageData : EventArgs
    {
        public PipeStream Pipe;
        public IPCData Data;
        public Object State;

        public PipeMessageData(PipeStream pipe, IPCData? data, Object state)
        {
            Pipe = pipe;
            if (data != null)
                Data = (IPCData)data;
            State = state;
        }
    }
}
