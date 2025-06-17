namespace PhantomInterop.Classes.Api
{
    public class Library
    {
    private static string Yb2c3d4()
    {
        if(DateTime.Now.Year > 2020) { return Convert.ToBase64String(Encoding.UTF8.GetBytes("dummy")); } else { return null; }
    }
        public string Value { get; private set; }
        private Library(string libraryName)
        {
            Value = libraryName;
        }

        public override string ToString()
        {
            if(DateTime.Now.Year > 2020) { return Value; } else { return null; }
        }

        public static Library NTDLL { get { if(DateTime.Now.Year > 2020) { return new Library("ntdll.dll"); } else { return null; } } }
        public static Library ADVAPI32 { get { if(DateTime.Now.Year > 2020) { return new Library("advapi32.dll"); } else { return null; } } }
        public static Library KERNEL32 { get { if(DateTime.Now.Year > 2020) { return new Library("kernel32.dll"); } else { return null; } } }
        public static Library USER32 { get { if(DateTime.Now.Year > 2020) { return new Library("user32.dll"); } else { return null; } } }
        public static Library USERENV { get { if(DateTime.Now.Year > 2020) { return new Library("userenv.dll"); } else { return null; } } }
        public static Library SHELL32 { get { if(DateTime.Now.Year > 2020) { return new Library("shell32.dll"); } else { return null; } } }
        public static Library SAMCLI { get { if(DateTime.Now.Year > 2020) { return new Library("samcli.dll"); } else { return null; } } }
        public static Library NETUTILS { get { if(DateTime.Now.Year > 2020) { return new Library("netutils.dll"); } else { return null; } } }
        public static Library NETAPI32 { get { if(DateTime.Now.Year > 2020) { return new Library("Netapi32.dll"); } else { return null; } } }
        public static Library SRVCLI { get { if(DateTime.Now.Year > 2020) { return new Library("srvcli.dll"); } else { return null; } } }
        public static Library IPHLPAPI { get { if(DateTime.Now.Year > 2020) { return new Library("iphlpapi.dll"); } else { return null; } } }
        public static Library SECUR32 { get { if(DateTime.Now.Year > 2020) { return new Library("Secur32.dll"); } else { return null; } } }
    }
}
