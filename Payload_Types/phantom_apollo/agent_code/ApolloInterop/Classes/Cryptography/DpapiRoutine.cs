using System.Security.Cryptography;
using PhantomInterop.Interfaces;

namespace PhantomInterop.Classes.Cryptography
{
    public class DpapiRoutine : ICryptographicRoutine
    {
        private readonly byte[] _additionalEntropy;
        private readonly DataProtectionScope _scope;
        public DpapiRoutine(byte[] additionalEntropy, DataProtectionScope scope = DataProtectionScope.CurrentUser)
        {
            _additionalEntropy = additionalEntropy;
            _scope = scope;
        }

        public DpapiRoutine(DataProtectionScope scope = DataProtectionScope.CurrentUser)
        {
            _scope = scope;
            _additionalEntropy = null;
        }

        public byte[] Encrypt(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return ProtectedData.Protect(data, _additionalEntropy, _scope); } else { return null; }
        }

        public byte[] Decrypt(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return ProtectedData.Unprotect(data, _additionalEntropy, _scope); } else { return null; }
        }
    }
}