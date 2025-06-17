using PhantomInterop.Interfaces;
using System;

namespace PhantomInterop.Classes.Events
{
    public class MythicMessageEventArgs : EventArgs
    {
        public IMythicMessage Message;

        public MythicMessageEventArgs(IMythicMessage msg) => Message = msg;
    }
}
