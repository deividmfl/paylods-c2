using PhantomInterop.Structs.MythicStructs;
using System.Net.Sockets;
using PhantomInterop.Classes;

namespace PhantomInterop.Interfaces
{
    public interface IRpfwdManager
    {
        bool Route(SocksDatagram dg);
        bool AddConnection(TcpClient client, int ServerID, int port, int debugLevel, Tasking task);
        bool Remove(int id);
    }
}
