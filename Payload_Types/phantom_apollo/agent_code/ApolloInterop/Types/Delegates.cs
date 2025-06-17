using PhantomInterop.Enums.PhantomEnums;

namespace PhantomInterop.Types
{
    namespace Delegates
    {
        public delegate bool OnResponse<T>(T message);
        public delegate bool DispatchMessage(byte[] data, MessageType mt);
    }
}
