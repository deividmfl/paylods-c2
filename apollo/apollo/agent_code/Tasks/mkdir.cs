#define COMMAND_NAME_UPPER

#if DEBUG
#define MKDIR
#endif

#if MKDIR

using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.IO;
using System.Runtime.Serialization;


namespace Tasks
{
    public class mkdir : Tasking
    {
        [DataContract]
        internal struct MkdirParameters
        {
            [DataMember(Name = "path")] public string Path;
        }
        
        public mkdir(IAgent agent, MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            MythicTaskResponse resp;
            MkdirParameters parameters = _dataSerializer.Deserialize<MkdirParameters>(_data.Parameters);
            
            if (System.IO.Directory.Exists(parameters.Path))
            {
                resp = CreateTaskResponse(
                    $"Directory {parameters.Path} already exists.",
                    true,
                    "error");
            }
            else
            {
                try
                {
                    DirectoryInfo info = Directory.CreateDirectory(parameters.Path);
                    FileInformation finfo = new FileInformation(info);
                    ICommandMessage[] artifacts = new ICommandMessage[2]
                    {
                        Artifact.FileCreate(info.FullName),
                        new FileBrowser(finfo),
                    };
                    resp = CreateTaskResponse(
                        $"Created {info.FullName}",
                        true,
                        "completed",
                        artifacts);
                }
                catch (Exception ex)
                {
                    resp = CreateTaskResponse(
                        $"Error creating {parameters.Path}: {ex.Message}",
                        true,
                        "error");
                }
            }

            
            
            
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }
    }
}

#endif