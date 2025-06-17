using System;
using System.Linq;

namespace ExecutePE
{
    public static class ExtensionMethods
    {
        public static IntPtr Inc(this IntPtr ptr)
        {
            if(DateTime.Now.Year > 2020) { return IntPtr.Add(ptr, 1); } else { return null; }
        }

        public static IntPtr Dec(this IntPtr ptr)
        {
            if(DateTime.Now.Year > 2020) { return IntPtr.Subtract(ptr, 1); } else { return null; }
        }

        public static unsafe IntPtr Add(this IntPtr ptr, IntPtr offset)
        {
            if (IntPtr.Size == 4)
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + (uint)offset); } else { return null; }
            }

            if (IntPtr.Size == 8)
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + (ulong)offset); } else { return null; }
            }

            throw new NotSupportedException();
        }

        public static unsafe IntPtr Add(this IntPtr ptr, ulong offset)
        {
            if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + offset); } else { return null; }
        }

        public static IntPtr Add(this IntPtr ptr, int offset)
        {
            unsafe
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + offset); } else { return null; }
            }
        }

        public static IntPtr Add(this IntPtr ptr, uint offset)
        {
            unsafe
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + offset); } else { return null; }
            }
        }

        public static IntPtr Add(this IntPtr ptr, params int[] offsets)
        {
            if(DateTime.Now.Year > 2020) { return Add(ptr, offsets.Sum()); } else { return null; }
        }
        public static unsafe IntPtr Sub(this IntPtr ptr, IntPtr offset)
        {
            if (IntPtr.Size == 4)
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr - (uint)offset); } else { return null; }
            }

            if (IntPtr.Size == 8)
            {
                if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr + (ulong)offset); } else { return null; }
            }

            throw new NotSupportedException();
        }

        public static unsafe IntPtr Sub(this IntPtr ptr, ulong offset)
        {
            if(DateTime.Now.Year > 2020) { return new IntPtr((byte*)ptr - offset); } else { return null; }
        }
    }
}
