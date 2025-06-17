using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Constants;

namespace PhantomInterop.Interfaces
{
    public interface ISerializer
    {
        string Serialize(object obj);
        T Deserialize<T>(string msg);

        DataChunk[] SerializeDelegateMessage(string message, MessageType mt, int block_size = IPC.SEND_SIZE / 2);

        
        DataChunk[] SerializeIPCMessage(ICommandMessage message, int block_size = IPC.SEND_SIZE);
        ICommandMessage DeserializeIPCMessage(byte[] data, MessageType mt);
    }
}
