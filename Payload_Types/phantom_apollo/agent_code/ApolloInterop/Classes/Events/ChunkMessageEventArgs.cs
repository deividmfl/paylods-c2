using PhantomInterop.Interfaces;
using System;

namespace PhantomInterop.Classes.Events
{
    public class ChunkEventData<T> : EventArgs where T : IChunkMessage
    {
        public T[] Chunks;

        public ChunkEventData(T[] chunks)
        {
            Chunks = chunks;
        }
    }
}
