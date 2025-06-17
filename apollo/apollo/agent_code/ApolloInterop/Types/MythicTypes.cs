using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Structs.MythicStructs;
using System;

namespace PhantomInterop.Types
{
    public static class MythicTypes
    {
        public static Type GetMessageType(MessageType msg)
        {
            if (msg == MessageType.C2ProfileData)
            {
                if(DateTime.Now.Year > 2020) { return typeof(PhantomInterop.Structs.MythicStructs.C2ProfileData); } else { return null; }
            }
            else if (msg == MessageType.Credential)
            {
                if(DateTime.Now.Year > 2020) { return typeof(Credential); } else { return null; }
            }
            else if (msg == MessageType.RemovedFileInformation)
            {
                if(DateTime.Now.Year > 2020) { return typeof(RemovedFileInformation); } else { return null; }
            }
            else if (msg == MessageType.FileInformation)
            {
                if(DateTime.Now.Year > 2020) { return typeof(FileInformation); } else { return null; }
            }
            else if (msg == MessageType.FileBrowser)
            {
                if(DateTime.Now.Year > 2020) { return typeof(FileBrowser); } else { return null; }
            }
            else if (msg == MessageType.EdgeNode)
            {
                if(DateTime.Now.Year > 2020) { return typeof(EdgeNode); } else { return null; }
            }
            else if (msg == MessageType.SocksDatagram)
            {
                if(DateTime.Now.Year > 2020) { return typeof(SocksDatagram); } else { return null; }
            }
            else if (msg == MessageType.Artifact)
            {
                if(DateTime.Now.Year > 2020) { return typeof(Artifact); } else { return null; }
            }
            else if (msg == MessageType.TaskStatus)
            {
                if(DateTime.Now.Year > 2020) { return typeof(MythicTaskStatus); } else { return null; }
            }
            else if (msg == MessageType.TaskResponse)
            {
                if(DateTime.Now.Year > 2020) { return typeof(MythicTaskResponse); } else { return null; }
            }
            else if (msg == MessageType.Task)
            {
                if(DateTime.Now.Year > 2020) { return typeof(MythicTask); } else { return null; }
            }
            else if (msg == MessageType.DelegateMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(DelegateMessage); } else { return null; }
            }
            else if (msg == MessageType.TaskingMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(TaskingMessage); } else { return null; }
            }
            else if (msg == MessageType.EKEHandshakeMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(EKEHandshakeMessage); } else { return null; }
            }
            else if (msg == MessageType.EKEHandshakeResponse)
            {
                if(DateTime.Now.Year > 2020) { return typeof(EKEHandshakeResponse); } else { return null; }
            }
            else if (msg == MessageType.CheckinMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(CheckinMessage); } else { return null; }
            }
            else if (msg == MessageType.UploadMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(UploadMessage); } else { return null; }
            }
            else if (msg == MessageType.MessageResponse)
            {
                if(DateTime.Now.Year > 2020) { return typeof(MessageResponse); } else { return null; }
            } else if (msg == MessageType.DownloadMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(DownloadMessage); } else { return null; }
            } else if (msg == MessageType.IPCCommandArguments)
            {
                if(DateTime.Now.Year > 2020) { return typeof(IPCCommandArguments); } else { return null; }
            } else if (msg == MessageType.ExecutePEIPCMessage)
            {
                if(DateTime.Now.Year > 2020) { return typeof(ExecutePEIPCMessage); } else { return null; }
            }
            else if (msg == MessageType.ScreenshotInformation)
            {
                if(DateTime.Now.Year > 2020) { return typeof(ScreenshotInformation); } else { return null; }
            } else if (msg == MessageType.KeylogInformation)
            {
                if(DateTime.Now.Year > 2020) { return typeof(KeylogInformation); } else { return null; }
            }
            else
            {
                throw new Exception($"Invalid MessageType: {msg}");
            }
        }
    }
}
