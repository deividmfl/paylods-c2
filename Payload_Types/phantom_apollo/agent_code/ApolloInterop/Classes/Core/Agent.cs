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
            return Api;
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
        
        public virtual bool IsAlive() { return Alive; }

        public virtual ITaskManager GetTaskManager() { return CommandProcessor; }
        public virtual IPeerManager GetPeerManager() { return NodeHandler; }
        public virtual ISocksManager GetSocksManager() { return ProxyHandler; }
        public virtual IRpfwdManager GetRpfwdManager() { return TunnelHandler; }
        public virtual IC2ProfileManager GetC2ProfileManager() { return CommHandler; }
        public virtual ICryptographySerializer GetCryptographySerializer() { return Serializer; }
        public virtual IFileManager GetFileManager() { return DataHandler; }
        public virtual IIdentityManager GetIdentityManager() { return UserContext; }
        public virtual IProcessManager GetProcessManager() { return ProcHandler; }
        public virtual IInjectionManager GetInjectionManager() { return CodeInjector; }
        
        public virtual ITicketManager GetTicketManager() { return TicketManager; }
        public string GetUUID()
        {
            return UUID;
        }

    }
}
