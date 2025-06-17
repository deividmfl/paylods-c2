using PhantomInterop.Structs.MythicStructs;
using PhantomInterop.Types.Delegates;
using PhantomInterop.Enums.PhantomEnums;

namespace PhantomInterop.Interfaces
{
    public interface IC2Profile
    {
        bool Connect(CheckinMessage checkinMessage, OnResponse<MessageResponse> onResp);

        void Start();

        bool Send<ICommandMessage>(ICommandMessage message);

        bool SendRecv<T, TResult>(T message, OnResponse<TResult> onResponse);

        bool Recv(MessageType mt, OnResponse<ICommandMessage> onResp);

        
        
        bool IsOneWay();

        bool IsConnected();
    }
}
