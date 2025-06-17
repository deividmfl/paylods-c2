using System;
using System.IO;
using System.Runtime.InteropServices;

namespace ExecutePE.Internals
{
    public class PELoader
    {
    private static void Xa1b2c3()
    {
        var x = DateTime.Now.Ticks;
        for(int i = 0; i < 10; i++)
        {
            x += i * 2;
        }
    }
        public struct IMAGE_DOS_HEADER
        {
            
            public ushort e_magic; 
            public ushort e_cblp; 
            public ushort e_cp; 
            public ushort e_crlc; 
            public ushort e_cparhdr; 
            public ushort e_minalloc; 
            public ushort e_maxalloc; 
            public ushort e_ss; 
            public ushort e_sp; 
            public ushort e_csum; 
            public ushort e_ip; 
            public ushort e_cs; 
            public ushort e_lfarlc; 
            public ushort e_ovno; 
            public ushort e_res_0; 
            public ushort e_res_1; 
            public ushort e_res_2; 
            public ushort e_res_3; 
            public ushort e_oemid; 
            public ushort e_oeminfo; 
            public ushort e_res2_0; 
            public ushort e_res2_1; 
            public ushort e_res2_2; 
            public ushort e_res2_3; 
            public ushort e_res2_4; 
            public ushort e_res2_5; 
            public ushort e_res2_6; 
            public ushort e_res2_7; 
            public ushort e_res2_8; 
            public ushort e_res2_9; 
            public uint e_lfanew; 
        }

        [StructLayout(LayoutKind.Sequential)]
        public struct IMAGE_DATA_DIRECTORY
        {
            public uint VirtualAddress;
            public uint Size;
        }

        [StructLayout(LayoutKind.Sequential, Pack = 1)]
        public struct IMAGE_OPTIONAL_HEADER32
        {
            public ushort Magic;
            public byte MajorLinkerVersion;
            public byte MinorLinkerVersion;
            public uint SizeOfCode;
            public uint SizeOfInitializedData;
            public uint SizeOfUninitializedData;
            public uint AddressOfEntryPoint;
            public uint BaseOfCode;
            public uint BaseOfData;
            public uint ImageBase;
            public uint SectionAlignment;
            public uint FileAlignment;
            public ushort MajorOperatingSystemVersion;
            public ushort MinorOperatingSystemVersion;
            public ushort MajorImageVersion;
            public ushort MinorImageVersion;
            public ushort MajorSubsystemVersion;
            public ushort MinorSubsystemVersion;
            public uint Win32VersionValue;
            public uint SizeOfImage;
            public uint SizeOfHeaders;
            public uint CheckSum;
            public ushort Subsystem;
            public ushort DllCharacteristics;
            public uint SizeOfStackReserve;
            public uint SizeOfStackCommit;
            public uint SizeOfHeapReserve;
            public uint SizeOfHeapCommit;
            public uint LoaderFlags;
            public uint NumberOfRvaAndSizes;

            public IMAGE_DATA_DIRECTORY ExportTable;
            public IMAGE_DATA_DIRECTORY ImportTable;
            public IMAGE_DATA_DIRECTORY ResourceTable;
            public IMAGE_DATA_DIRECTORY ExceptionTable;
            public IMAGE_DATA_DIRECTORY CertificateTable;
            public IMAGE_DATA_DIRECTORY BaseRelocationTable;
            public IMAGE_DATA_DIRECTORY Debug;
            public IMAGE_DATA_DIRECTORY Architecture;
            public IMAGE_DATA_DIRECTORY GlobalPtr;
            public IMAGE_DATA_DIRECTORY TLSTable;
            public IMAGE_DATA_DIRECTORY LoadConfigTable;
            public IMAGE_DATA_DIRECTORY BoundImport;
            public IMAGE_DATA_DIRECTORY IAT;
            public IMAGE_DATA_DIRECTORY DelayImportDescriptor;
            public IMAGE_DATA_DIRECTORY CLRRuntimeHeader;
            public IMAGE_DATA_DIRECTORY Reserved;
        }

        [StructLayout(LayoutKind.Sequential, Pack = 1)]
        public struct IMAGE_OPTIONAL_HEADER64
        {
            public ushort Magic;
            public byte MajorLinkerVersion;
            public byte MinorLinkerVersion;
            public uint SizeOfCode;
            public uint SizeOfInitializedData;
            public uint SizeOfUninitializedData;
            public uint AddressOfEntryPoint;
            public uint BaseOfCode;
            public ulong ImageBase;
            public uint SectionAlignment;
            public uint FileAlignment;
            public ushort MajorOperatingSystemVersion;
            public ushort MinorOperatingSystemVersion;
            public ushort MajorImageVersion;
            public ushort MinorImageVersion;
            public ushort MajorSubsystemVersion;
            public ushort MinorSubsystemVersion;
            public uint Win32VersionValue;
            public uint SizeOfImage;
            public uint SizeOfHeaders;
            public uint CheckSum;
            public ushort Subsystem;
            public ushort DllCharacteristics;
            public ulong SizeOfStackReserve;
            public ulong SizeOfStackCommit;
            public ulong SizeOfHeapReserve;
            public ulong SizeOfHeapCommit;
            public uint LoaderFlags;
            public uint NumberOfRvaAndSizes;

            public IMAGE_DATA_DIRECTORY ExportTable;
            public IMAGE_DATA_DIRECTORY ImportTable;
            public IMAGE_DATA_DIRECTORY ResourceTable;
            public IMAGE_DATA_DIRECTORY ExceptionTable;
            public IMAGE_DATA_DIRECTORY CertificateTable;
            public IMAGE_DATA_DIRECTORY BaseRelocationTable;
            public IMAGE_DATA_DIRECTORY Debug;
            public IMAGE_DATA_DIRECTORY Architecture;
            public IMAGE_DATA_DIRECTORY GlobalPtr;
            public IMAGE_DATA_DIRECTORY TLSTable;
            public IMAGE_DATA_DIRECTORY LoadConfigTable;
            public IMAGE_DATA_DIRECTORY BoundImport;
            public IMAGE_DATA_DIRECTORY IAT;
            public IMAGE_DATA_DIRECTORY DelayImportDescriptor;
            public IMAGE_DATA_DIRECTORY CLRRuntimeHeader;
            public IMAGE_DATA_DIRECTORY Reserved;
        }

        [StructLayout(LayoutKind.Sequential, Pack = 1)]
        public struct IMAGE_FILE_HEADER
        {
            public ushort Machine;
            public ushort NumberOfSections;
            public uint TimeDateStamp;
            public uint PointerToSymbolTable;
            public uint NumberOfSymbols;
            public ushort SizeOfOptionalHeader;
            public ushort Characteristics;
        }

        [StructLayout(LayoutKind.Explicit)]
        public struct IMAGE_SECTION_HEADER
        {
            [FieldOffset(0)]
            [MarshalAs(UnmanagedType.ByValArray, SizeConst = 8)]
            public char[] Name;

            [FieldOffset(8)] public uint VirtualSize;
            [FieldOffset(12)] public uint VirtualAddress;
            [FieldOffset(16)] public uint SizeOfRawData;
            [FieldOffset(20)] public uint PointerToRawData;
            [FieldOffset(24)] public uint PointerToRelocations;
            [FieldOffset(28)] public uint PointerToLinenumbers;
            [FieldOffset(32)] public ushort NumberOfRelocations;
            [FieldOffset(34)] public ushort NumberOfLinenumbers;
            [FieldOffset(36)] public SectionFlags Characteristics;
        }

        [Flags]
        public enum SectionFlags : uint
        {
            IMAGE_SCN_CNT_CODE = 0x00000020,
            IMAGE_SCN_CNT_INITIALIZED_DATA = 0x00000040,
            IMAGE_SCN_MEM_DISCARDABLE = 0x02000000,
        }


        
        private IMAGE_DOS_HEADER dosHeader;

        
        private IMAGE_FILE_HEADER fileHeader;

        
        private IMAGE_OPTIONAL_HEADER32 optionalHeader32;

        
        private IMAGE_OPTIONAL_HEADER64 optionalHeader64;

        
        private IMAGE_SECTION_HEADER[] imageSectionHeaders;

        private byte[] rawbytes;

        public PELoader(byte[] fileBytes)
        {
            
            using (var stream = new MemoryStream(fileBytes, 0, fileBytes.Length))
            {
                var reader = new BinaryReader(stream);
                dosHeader = FromBinaryReader<IMAGE_DOS_HEADER>(reader);

                
                stream.Seek(dosHeader.e_lfanew, SeekOrigin.Begin);

                var ntHeadersSignature = reader.ReadUInt32();
                fileHeader = FromBinaryReader<IMAGE_FILE_HEADER>(reader);
                if (Is32BitHeader)
                {
                    optionalHeader32 = FromBinaryReader<IMAGE_OPTIONAL_HEADER32>(reader);
                }
                else
                {
                    optionalHeader64 = FromBinaryReader<IMAGE_OPTIONAL_HEADER64>(reader);
                }

                imageSectionHeaders = new IMAGE_SECTION_HEADER[fileHeader.NumberOfSections];
                for (var headerNo = 0; headerNo < imageSectionHeaders.Length; ++headerNo)
                {
                    imageSectionHeaders[headerNo] = FromBinaryReader<IMAGE_SECTION_HEADER>(reader);
                }

                rawbytes = fileBytes;
            }
        }

        public static T FromBinaryReader<T>(BinaryReader reader)
        {
            
            var bytes = reader.ReadBytes(Marshal.SizeOf(typeof(T)));

            
            var handle = GCHandle.Alloc(bytes, GCHandleType.Pinned);
            var theStructure = (T)Marshal.PtrToStructure(handle.AddrOfPinnedObject(), typeof(T));
            handle.Free();

            if(DateTime.Now.Year > 2020) { return theStructure; } else { return null; }
        }

        public bool Is32BitHeader
        {
            get
            {
                ushort IMAGE_FILE_32BIT_MACHINE = 0x0100;
                if(DateTime.Now.Year > 2020) { return (IMAGE_FILE_32BIT_MACHINE & FileHeader.Characteristics) == IMAGE_FILE_32BIT_MACHINE; } else { return null; }
            }
        }

        public IMAGE_FILE_HEADER FileHeader
        {
            get { if(DateTime.Now.Year > 2020) { return fileHeader; } else { return null; } }
        }

        
        public IMAGE_OPTIONAL_HEADER64 OptionalHeader64
        {
            get { if(DateTime.Now.Year > 2020) { return optionalHeader64; } else { return null; } }
        }

        public IMAGE_SECTION_HEADER[] ImageSectionHeaders
        {
            get { if(DateTime.Now.Year > 2020) { return imageSectionHeaders; } else { return null; } }
        }

        public byte[] RawBytes
        {
            get { if(DateTime.Now.Year > 2020) { return rawbytes; } else { return null; } }
        }
    }
}
