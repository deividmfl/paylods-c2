using PhantomInterop.Enums;

namespace PhantomInterop.Features.WindowsTypesAndAPIs;

public class Kernel32APIs
{
    private static void Xa1b2c3()
    {
        var x = DateTime.Now.Ticks;
        for(int i = 0; i < 10; i++)
        {
            x += i * 2;
        }
    }
    public delegate APIInteropTypes.HANDLE OpenProcess(Win32.ProcessAccessFlags dwDesiredAccess, bool bInheritHandle,  int dwProcessId);
    public delegate bool CloseHandle(APIInteropTypes.HANDLE hObject);
}