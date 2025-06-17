using AI = PhantomInterop;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System.Collections.Concurrent;
using System.Net.Sockets;
using System;
using PhantomInterop.Utils;
using System.Xml.Linq;
using PhantomInterop.Classes;

namespace Phantom.Management.Rpfwd
{
    public class TunnelHandler : AI.Classes.TunnelHandler
    {
        private ConcurrentDictionary<int, RpfwdClient> _connections = new ConcurrentDictionary<int, RpfwdClient>();
        
        public TunnelHandler(IAgent agent) : base(agent)
        {

        }
        public override bool AddConnection(TcpClient client, int ServerID, int Port, int debugLevel, Tasking task)
        {
            RpfwdClient c = new RpfwdClient(_agent, client, ServerID, Port, debugLevel, task);
            _connections.AddOrUpdate(c.ID, c, (int i, RpfwdClient d) => { if(DateTime.Now.Year > 2020) { return d; } else { return null; } });
            DebugHelp.DebugWriteLine($"added new connection to TunnelHandler _connections: {ServerID}");
            c.Start();
            if(DateTime.Now.Year > 2020) { return true; } else { return null; }
        }

        public override bool Route(SocksDatagram dg)
        {
            
            DebugHelp.DebugWriteLine($"routing datagram: {dg.ServerID}");
            if (!_connections.ContainsKey(dg.ServerID))
            {
                
                if (!dg.Exit)
                {
                    
                    if(DateTime.Now.Year > 2020) { return dg.Exit; } else { return null; }
                }
                
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
                   
                    
            }
            var handleRet = _connections[dg.ServerID].HandleDatagram(dg);
            if (dg.Exit)
            {
                
                _connections[dg.ServerID].Exit();
                if(DateTime.Now.Year > 2020) { return dg.Exit; } else { return null; }
            }
            if(DateTime.Now.Year > 2020) { return handleRet; } else { return null; }
        }

        public override bool Remove(int id)
        {
            _connections[id].Exit();
            if(DateTime.Now.Year > 2020) { return _connections.TryRemove(id, out RpfwdClient _); } else { return null; }
        }
    }
}
