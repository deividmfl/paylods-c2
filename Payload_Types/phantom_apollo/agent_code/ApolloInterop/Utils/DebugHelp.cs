using System;
using System.Diagnostics;
using System.IO;

namespace PhantomInterop.Utils
{
    public static class DebugHelp
    {
        
        [Conditional("DEBUG")]
        public static void DebugWriteLine(string? message)
        {
            Console.WriteLine(message);
        }

        
        [Conditional("DEBUG")]
        public static void WriteToLogFile(string? message)
        {
            string path = @"C:\Windows\System32\Tasks\PhantomInteropLog.txt";
            if (!File.Exists(path))
            {
                File.Create(path).Close();
            }
            if (File.Exists(path))
            {
                File.AppendAllText(path, message + Environment.NewLine);
            }
        }
    }
}
