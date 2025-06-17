using PhantomInterop.Classes.Events;
using PhantomInterop.Interfaces;
using System;

namespace PhantomInterop.Classes.Core
{
    public class ChunkStore<T> where T : IChunkMessage
    {
        private T[] _messages = null;
        private object _lock = new object();
        private int _currentCount = 0;

        public event EventHandler<ChunkEventData<T>> ChunkAdd;
        public event EventHandler<ChunkEventData<T>> MessageComplete;
        public void OnMessageComplete() => MessageComplete?.Invoke(this, new ChunkEventData<T>(_messages));
        public void AddMessage(T d)
        {
            lock(_lock)
            {
                if (_messages == null)
                {
                    _messages = new T[d.GetTotalChunks()];
                }
                _messages[d.GetChunkNumber()-1] = d;
                _currentCount += 1;
            }
            if (_currentCount == d.GetTotalChunks())
            {
                OnMessageComplete();
            } else
            {
                ChunkAdd?.Invoke(this, new ChunkEventData<T>(new T[1] { d }));
            }
        }
    }
}
