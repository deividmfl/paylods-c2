#define COMMAND_NAME_UPPER

#if DEBUG
#define GET_INJECTION_TECHNIQUES
#endif

#if GET_INJECTION_TECHNIQUES

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.Collections.Generic;
using System.Runtime.Serialization;

namespace Tasks
{
    public class get_injection_techniques : Tasking
    {
        [DataContract]
        internal struct InjectionTechniqueResult
        {
            [DataMember(Name = "name")]
            public string Name;
            [DataMember(Name = "is_current")]
            public bool IsCurrent;
        }
        public get_injection_techniques(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }
        

        public override void Start()
        {
            MythicTaskResponse resp;
            string[] techniques = _agent.GetInjectionManager().GetTechniques();
            Type cur = _agent.GetInjectionManager().GetCurrentTechnique();
            List<InjectionTechniqueResult> results = new List<InjectionTechniqueResult>();
            foreach (string t in techniques)
            {
                results.Add(new InjectionTechniqueResult
                {
                    Name = t,
                    IsCurrent = t == cur.Name
                });
            }

            resp = CreateTaskResponse(
                _dataSerializer.Serialize(results.ToArray()), true);
            
            
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}

#endif