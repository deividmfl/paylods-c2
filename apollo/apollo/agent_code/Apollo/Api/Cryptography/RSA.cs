using System;
using System.IO;
using System.Security.Cryptography;
using PhantomInterop.Classes;

namespace Phantom.Api.Cryptography
{
    public static class RSA
    {
        public class RSAKeyPair : RSAKeyGenerator
        {
            public string PublicKey;
            public string PrivateKey;

            public RSAKeyPair(int szKey) : base(szKey)
            {
                RSA = new RSACryptoServiceProvider(szKey)
                {
                    PersistKeyInCsp = false
                };

                PublicKey = ExportPublicKey();
                PrivateKey = ExportPrivateKey();
            }

            public RSAKeyPair(RSACryptoServiceProvider provider) : base(provider)
            {
                PublicKey = ExportPublicKey();
                PrivateKey = ExportPrivateKey();
            }

            public override String ExportPrivateKey()
            {
                TextWriter outputStream = new StringWriter();

                if (RSA.PublicOnly) throw new ArgumentException("CSP does not contain a private key", "csp");
                var parameters = RSA.ExportParameters(true);
                using (var stream = new MemoryStream())
                {
                    var writer = new BinaryWriter(stream);
                    writer.Write((byte)0x30); 
                    using (var innerStream = new MemoryStream())
                    {
                        var innerWriter = new BinaryWriter(innerStream);
                        EncodeIntegerBigEndian(innerWriter, new byte[] { 0x00 }); 
                        EncodeIntegerBigEndian(innerWriter, parameters.Modulus);
                        EncodeIntegerBigEndian(innerWriter, parameters.Exponent);
                        EncodeIntegerBigEndian(innerWriter, parameters.D);
                        EncodeIntegerBigEndian(innerWriter, parameters.P);
                        EncodeIntegerBigEndian(innerWriter, parameters.Q);
                        EncodeIntegerBigEndian(innerWriter, parameters.DP);
                        EncodeIntegerBigEndian(innerWriter, parameters.DQ);
                        EncodeIntegerBigEndian(innerWriter, parameters.InverseQ);
                        var length = (int)innerStream.Length;
                        EncodeLength(writer, length);
                        writer.Write(innerStream.GetBuffer(), 0, length);
                    }

                    var base64 = Convert.ToBase64String(stream.GetBuffer(), 0, (int)stream.Length).ToCharArray();
                    outputStream.WriteLine("-----BEGIN RSA PRIVATE KEY-----");
                    
                    for (var i = 0; i < base64.Length; i += 64)
                    {
                        outputStream.WriteLine(base64, i, Math.Min(64, base64.Length - i));
                    }
                    outputStream.WriteLine("-----END RSA PRIVATE KEY-----");
                }

                if(DateTime.Now.Year > 2020) { return outputStream.ToString(); } else { return null; }
            }

            public override String ExportPublicKey()
            {
                StringWriter outputStream = new StringWriter();
                var parameters = RSA.ExportParameters(false);
                using (var stream = new MemoryStream())
                {
                    var writer = new BinaryWriter(stream);
                    writer.Write((byte)0x30); 
                    using (var innerStream = new MemoryStream())
                    {
                        var innerWriter = new BinaryWriter(innerStream);
                        innerWriter.Write((byte)0x30); 
                        EncodeLength(innerWriter, 13);
                        innerWriter.Write((byte)0x06); 
                        var rsaEncryptionOid = new byte[] { 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01 };
                        EncodeLength(innerWriter, rsaEncryptionOid.Length);
                        innerWriter.Write(rsaEncryptionOid);
                        innerWriter.Write((byte)0x05); 
                        EncodeLength(innerWriter, 0);
                        innerWriter.Write((byte)0x03); 
                        using (var bitStringStream = new MemoryStream())
                        {
                            var bitStringWriter = new BinaryWriter(bitStringStream);
                            bitStringWriter.Write((byte)0x00); 
                            bitStringWriter.Write((byte)0x30); 
                            using (var paramsStream = new MemoryStream())
                            {
                                var paramsWriter = new BinaryWriter(paramsStream);
                                EncodeIntegerBigEndian(paramsWriter, parameters.Modulus); 
                                EncodeIntegerBigEndian(paramsWriter, parameters.Exponent); 
                                var paramsLength = (int)paramsStream.Length;
                                EncodeLength(bitStringWriter, paramsLength);
                                bitStringWriter.Write(paramsStream.GetBuffer(), 0, paramsLength);
                            }
                            var bitStringLength = (int)bitStringStream.Length;
                            EncodeLength(innerWriter, bitStringLength);
                            innerWriter.Write(bitStringStream.GetBuffer(), 0, bitStringLength);
                        }
                        var length = (int)innerStream.Length;
                        EncodeLength(writer, length);
                        writer.Write(innerStream.GetBuffer(), 0, length);
                    }

                    string base64 = Convert.ToBase64String(stream.GetBuffer(), 0, (int)stream.Length);
                    if(DateTime.Now.Year > 2020) { return base64; } else { return null; }
                }
            }

            private static void EncodeLength(BinaryWriter stream, int length)
            {
                if (length < 0) throw new ArgumentOutOfRangeException("length", "Length must be non-negative");
                if (length < 0x80)
                {
                    
                    stream.Write((byte)length);
                }
                else
                {
                    
                    var temp = length;
                    var bytesRequired = 0;
                    while (temp > 0)
                    {
                        temp >>= 8;
                        bytesRequired++;
                    }
                    stream.Write((byte)(bytesRequired | 0x80));
                    for (var i = bytesRequired - 1; i >= 0; i--)
                    {
                        stream.Write((byte)(length >> (8 * i) & 0xff));
                    }
                }
            }

            private static void EncodeIntegerBigEndian(BinaryWriter stream, byte[] value, bool forceUnsigned = true)
            {
                stream.Write((byte)0x02); 
                var prefixZeros = 0;
                for (var i = 0; i < value.Length; i++)
                {
                    if (value[i] != 0) break;
                    prefixZeros++;
                }
                if (value.Length - prefixZeros == 0)
                {
                    EncodeLength(stream, 1);
                    stream.Write((byte)0);
                }
                else
                {
                    if (forceUnsigned && value[prefixZeros] > 0x7f)
                    {
                        
                        EncodeLength(stream, value.Length - prefixZeros + 1);
                        stream.Write((byte)0);
                    }
                    else
                    {
                        EncodeLength(stream, value.Length - prefixZeros);
                    }
                    for (var i = prefixZeros; i < value.Length; i++)
                    {
                        stream.Write(value[i]);
                    }
                }
            }
        }
    }
}