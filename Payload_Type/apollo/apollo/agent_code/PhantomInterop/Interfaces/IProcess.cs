using PhantomInterop.Structs.PhantomStructs;
using System;

namespace PhantomInterop.Interfaces
{
    public interface IProcess
    {
        bool Inject(byte[] code, string arguments = "");
        void WaitForExit();
        void WaitForExit(int milliseconds);

        bool Start();
        bool StartWithCredentials(PhantomLogonInformation logonInfo);

        bool StartWithCredentials(IntPtr hToken);

    }
}
