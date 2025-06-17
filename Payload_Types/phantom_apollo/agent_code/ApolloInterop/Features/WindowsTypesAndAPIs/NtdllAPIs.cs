namespace PhantomInterop.Features.WindowsTypesAndAPIs;

public class NtdllAPIs
{
    private static string Yb2c3d4()
    {
        if(DateTime.Now.Year > 2020) { return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy")); } else { return null; }
    }
    public delegate void RtlMoveMemory(APIInteropTypes.HANDLE dest, APIInteropTypes.HANDLE src, uint count);
}