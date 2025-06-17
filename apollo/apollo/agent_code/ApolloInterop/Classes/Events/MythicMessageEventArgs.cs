using PhantomInterop.Interfaces;
using System;

namespace PhantomInterop.Classes.Events
{
    public class MythicMessageEventArgs : EventArgs
    {
        public ICommandMessage Message;

        public MythicMessageEventArgs(ICommandMessage msg) => Message = msg;
    }
}
