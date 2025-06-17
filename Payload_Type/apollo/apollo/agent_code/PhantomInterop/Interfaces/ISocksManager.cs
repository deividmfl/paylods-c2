using PhantomInterop.Structs.MythicStructs;

namespace PhantomInterop.Interfaces
{
    public interface ISocksManager
    {
        bool Route(SocksDatagram dg);

        bool Remove(int id);
    }
}
