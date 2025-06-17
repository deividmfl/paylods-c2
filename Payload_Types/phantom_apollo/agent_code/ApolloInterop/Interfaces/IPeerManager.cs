using PhantomInterop.Classes.P2P;
using PhantomInterop.Structs.MythicStructs;
namespace PhantomInterop.Interfaces
{
    public interface IPeerManager
    {
        Peer AddPeer(PeerInformation info);
        bool Remove(string uuid);
        bool Remove(IPeer peer);
        bool Route(DelegateMessage msg);
    }
}
