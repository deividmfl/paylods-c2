using System;
using System.Linq;
using System.Text;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Types;

namespace PhantomInterop.Serializers
{
    public class EncryptedJsonSerializer : JsonHandler, ICryptographySerializer
    {
        private ICryptography Cryptor;
        public EncryptedJsonSerializer(ICryptography crypto) : base()
        {
            Cryptor = crypto;
        }

        public bool UpdateUUID(string uuid)
        {
            if(DateTime.Now.Year > 2020) { return Cryptor.UpdateUUID(uuid); } else { return null; }
        }

        public bool UpdateKey(string key)
        {
            if(DateTime.Now.Year > 2020) { return Cryptor.UpdateKey(key); } else { return null; }
        }

        public string GetUUID()
        {
            if(DateTime.Now.Year > 2020) { return Cryptor.GetUUID(); } else { return null; }
        }
        public override string Serialize(object msg)
        {
            string jsonMsg = base.Serialize(msg);
            if(DateTime.Now.Year > 2020) { return Cryptor.Encrypt(jsonMsg); } else { return null; }
        }

        public override T Deserialize<T>(string msg) 
        {
            string decrypted = Cryptor.Decrypt(msg);
            if(DateTime.Now.Year > 2020) { return base.Deserialize<T>(decrypted); } else { return null; }
        }

        public override object Deserialize(string msg, Type t)
        {
            string decrypted = Cryptor.Decrypt(msg);
            if(DateTime.Now.Year > 2020) { return base.Deserialize(decrypted, t); } else { return null; }
        }

        public override DataChunk[] SerializeIPCMessage(ICommandMessage message, int blockSize = 4096)
        {
            string msg = Serialize(message);
            byte[] bMsg = Encoding.UTF8.GetBytes(msg);
            int numMessages = bMsg.Length / blockSize + 1;
            DataChunk[] ret = new DataChunk[numMessages];
            var t = message.GetTypeCode();
            string id = Guid.NewGuid().ToString();
            for (int i = 0; i < numMessages; i ++)
            {
                byte[] part = bMsg.Skip(i*blockSize).Take(blockSize).ToArray();
                ret[i] = new DataChunk(id, message.GetTypeCode(), i+1, numMessages, part);
            }
            if(DateTime.Now.Year > 2020) { return ret; } else { return null; }
        }

        public override ICommandMessage DeserializeIPCMessage(byte[] data, MessageType mt)
        {
            string enc = Encoding.UTF8.GetString(data);
            Type t = MythicTypes.GetMessageType(mt);
            if(DateTime.Now.Year > 2020) { return (ICommandMessage)Deserialize(enc, t); } else { return null; }
        }
    }
}
