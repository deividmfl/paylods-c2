using PhantomInterop.Classes;
using System;
using PhantomInterop.Classes.Api;

namespace PhantomInterop.Interfaces
{
    public interface IApi
    {
        T GetLibraryFunction<T>(Library library, string functionName, bool canLoadFromDisk = true, bool resolveForwards = true) where T : Delegate;
        T GetLibraryFunction<T>(Library library, short ordinal, bool canLoadFromDisk = true, bool resolveForwards = true) where T : Delegate;
        T GetLibraryFunction<T>(Library library, string functionHash, long key, bool canLoadFromDisk=true, bool resolveForwards = true) where T : Delegate;

        string NewUUID();

        RSAKeyGenerator NewRSAKeyPair(int szKey);

        
        ICryptographySerializer NewEncryptedJsonSerializer(string uuid, Type cryptoType, string key = "");
        
    }
}
