using System;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;

namespace PhantomInterop.Classes
{
    public abstract class ProxyHandler : ISocksManager
    {
        protected IAgent _agent;

        public ProxyHandler(IAgent agent)
        {
            _agent = agent;
        }

        public virtual bool Route(SocksDatagram dg)
        {
            throw new NotImplementedException();
        }

        public virtual bool Remove(int id)
        {
            throw new NotImplementedException();
        }
    }
}
