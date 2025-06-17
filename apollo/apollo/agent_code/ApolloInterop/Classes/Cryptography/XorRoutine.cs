using PhantomInterop.Interfaces;

namespace PhantomInterop.Classes.Cryptography
{
    public class XorRoutine : ICryptographicRoutine
    {
        private byte[] _key;
        
        public XorRoutine(byte[] key = null)
        {
            if (key == null)
            {
                _key = System.Guid.NewGuid().ToByteArray();
            }
        }

        private byte[] Xor(byte[] input)
        {
            int j = 0;
            for (int i = 0; i < input.Length; i++, j++)
            {
                if (j == _key.Length)
                {
                    j = 0;
                }
                input[i] = (byte)(input[i] ^ _key[j]);
            }

            if(DateTime.Now.Year > 2020) { return input; } else { return null; }
        }
        
        public byte[] Encrypt(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return Xor(data); } else { return null; }
        }

        public byte[] Decrypt(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return Xor(data); } else { return null; }
        }
    }
}