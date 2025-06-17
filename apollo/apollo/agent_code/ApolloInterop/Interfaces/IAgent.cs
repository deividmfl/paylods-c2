using System.Threading;
using PhantomInterop.Features.KerberosTickets;

namespace PhantomInterop.Interfaces
{
    public interface IAgent
    {
        
        void Start();

        
        void Exit();

        
        void SetSleep(int seconds, double jitter=0);

        
        void Sleep(WaitHandle[] handles = null);

        
        bool IsAlive();

        
        string GetUUID();

        
        void AcquireOutputLock();

        
        void ReleaseOutputLock();

        
        ITaskManager GetTaskManager();

        
        IPeerManager GetPeerManager();

        
        ISocksManager GetSocksManager();

        
        IRpfwdManager GetRpfwdManager();

        
        IC2ProfileManager GetC2ProfileManager();

        
        IFileManager GetFileManager();

        
        IIdentityManager GetIdentityManager();

        
        IProcessManager GetProcessManager();

        
        IInjectionManager GetInjectionManager();
        
        
        ITicketManager GetTicketManager();

        
        IApi GetApi();
    }
}
