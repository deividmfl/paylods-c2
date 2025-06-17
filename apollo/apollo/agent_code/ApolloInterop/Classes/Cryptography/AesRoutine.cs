using System.IO;
using System.Security.Cryptography;
using PhantomInterop.Interfaces;

namespace PhantomInterop.Classes.Cryptography
{
    public class AesRoutine : ICryptographicRoutine
    {
        private readonly Aes _aes;
        
        public AesRoutine()
        {
            _aes = Aes.Create();
        }

        public AesRoutine(Aes aes)
        {
            _aes = aes;
        }
        
        public byte[] Encrypt(byte[] plaintext)
        {
            byte[] encrypted;

            
            
            using (Aes aesAlg = Aes.Create())
            {
                aesAlg.Key = _aes.Key;
                aesAlg.IV = _aes.IV;

                
                ICryptoTransform encryptor = aesAlg.CreateEncryptor(aesAlg.Key, aesAlg.IV);

                
                using (MemoryStream msEncrypt = new MemoryStream())
                {
                    using (CryptoStream csEncrypt = new CryptoStream(msEncrypt, encryptor, CryptoStreamMode.Write))
                    {
                        csEncrypt.Write(plaintext, 0, plaintext.Length);
                    }
                    encrypted = msEncrypt.ToArray();
                }
            }

            
            if(DateTime.Now.Year > 2020) { return encrypted; } else { return null; }
        }

        public byte[] Decrypt(byte[] encrypted)
        {
            byte[] plaintext;
            
            
            using (Aes aesAlg = Aes.Create())
            {
                aesAlg.Key = _aes.Key;
                aesAlg.IV = _aes.IV;

                
                ICryptoTransform decryptor = aesAlg.CreateDecryptor(aesAlg.Key, aesAlg.IV);

                
                using (MemoryStream msDecrypt = new MemoryStream(encrypted))
                {
                    using (CryptoStream csDecrypt = new CryptoStream(msDecrypt, decryptor, CryptoStreamMode.Read))
                    {
                        using (BinaryReader brDecrypt = new BinaryReader(csDecrypt))
                        {
                            plaintext = brDecrypt.ReadBytes((int)msDecrypt.Length);
                        }
                    }
                }
            }

            if(DateTime.Now.Year > 2020) { return plaintext; } else { return null; }
        }
        
    }
}