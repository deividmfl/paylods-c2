﻿using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Utils;
using System;
using System.Collections.Concurrent;
using System.IO.Pipes;
using System.Security.AccessControl;
using System.Security.Principal;
using System.Threading;
using System.Threading.Tasks;

namespace PhantomInterop.Classes
{
    public class AsyncNamedPipeServer
    {
        private bool _running = true;

        private readonly string _pipeName;
        private readonly PipeSecurity _pipeSecurity;
        private readonly int _BUF_IN;
        private readonly int _BUF_OUT;
        private readonly int _maxInstances;

        private ConcurrentDictionary<PipeStream, IPCData> _connections = new();

        public event EventHandler<PipeMessageData> ConnectionEstablished;
        public event EventHandler<PipeMessageData> MessageReceived;
        public event EventHandler<PipeMessageData> Disconnect;

        public AsyncNamedPipeServer(string pipename, PipeSecurity ps = null, int instances=1, int BUF_IN=4096, int BUF_OUT=4096)
        {
            _pipeName = pipename;
            _BUF_IN = BUF_IN;
            _BUF_OUT = BUF_OUT;
            _maxInstances = instances;
            if (ps == null)
            {
                _pipeSecurity = new PipeSecurity();
                PipeAccessRule multipleInstances = new PipeAccessRule(WindowsIdentity.GetCurrent().Name, PipeAccessRights.CreateNewInstance, AccessControlType.Allow);
                PipeAccessRule everyoneAllowedRule = new PipeAccessRule(new SecurityIdentifier(WellKnownSidType.WorldSid, null), PipeAccessRights.ReadWrite, AccessControlType.Allow);
                PipeAccessRule networkAllowRule = new PipeAccessRule(new SecurityIdentifier(WellKnownSidType.NetworkSid, null), PipeAccessRights.ReadWrite, AccessControlType.Allow);
                _pipeSecurity.AddAccessRule(multipleInstances);
                _pipeSecurity.AddAccessRule(everyoneAllowedRule);
                _pipeSecurity.AddAccessRule(networkAllowRule);
            }

            for(int i = 0; i < _maxInstances; i++)
            {
                CreateServerPipe();
            }
        }

        public void Stop()
        {
            _running = false;
            foreach (var pipe in _connections.Keys)
            {
                pipe.Close();
            }
            while(true)
            {
                int count = _connections.Count;
                if (count == 0)
                    break;
                System.Threading.Thread.Sleep(5);
            }
        }

        private void CreateServerPipe()
        {
            DebugHelp.DebugWriteLine($"Creating Named Pipe: {_pipeName}");
            NamedPipeServerStream pipe = new NamedPipeServerStream(
                _pipeName,
                PipeDirection.InOut,
                -1,
                PipeTransmissionMode.Message,
                PipeOptions.Asynchronous | PipeOptions.WriteThrough,
                _BUF_IN,
                _BUF_OUT,
                _pipeSecurity
            );
            //NamedPipeServerStream pipe = new NamedPipeServerStream(_pipeName, PipeDirection.InOut, -1, PipeTransmissionMode.Message, PipeOptions.Asynchronous);
            // wait for client to connect async
            try
            {
                pipe.BeginWaitForConnection(OnClientConnected, pipe);
            }catch(Exception ex)
            {

            }

        }

        private bool IsPersistentConnectionAsync(NamedPipeServerStream pipeServer)
        {
            // Try to complete the read task within the timeout
            var completedTask = Task.WhenAny(
                    Task.Delay(500)
            );
            return pipeServer.IsConnected;
        }


        private void OnConnect(PipeMessageData args)
        {
            ConnectionEstablished?.Invoke(this, args);
        }

        private void OnMessageReceived(PipeMessageData args)
        {
            MessageReceived?.Invoke(this, args);
        }

        private void OnDisconnect(PipeMessageData args)
        {
            DebugHelp.DebugWriteLine($"Client disconnected from Named Pipe: {_pipeName}");
            Disconnect?.Invoke(this, args);
        }

        private void OnClientConnected(IAsyncResult result)
        {
            // complete connection
            NamedPipeServerStream pipe = (NamedPipeServerStream)result.AsyncState;
            pipe.EndWaitForConnection(result);
            DebugHelp.DebugWriteLine($"Client connected to Named Pipe: {_pipeName}");
            if(!IsPersistentConnectionAsync(pipe))
            {
                pipe.Close();
            }
            // create client pipe structure
            IPCData pd = new IPCData()
            {
                Pipe = pipe,
                State = null,
                Data = new byte[_BUF_IN],
            };
            
            // Add to connection list
            if (_running && _connections.TryAdd(pipe, pd))
            {
                // Prep the next connection
                CreateServerPipe();
                OnConnect(new PipeMessageData(pipe, null, this));
                BeginRead(pd);
            } else
            {
                pipe.Close();
            }
        }

        private void BeginRead(IPCData pd)
        {
            bool isConnected = pd.Pipe.IsConnected;
            if (isConnected)
            {
                try
                {
                    pd.Pipe.BeginRead(pd.Data, 0, pd.Data.Length, ProcessReceivedMessage, pd);
                } catch (Exception ex)
                {
                    isConnected = false;
                }
            }

            if (!isConnected)
            {
                pd.Pipe.Close();
                OnDisconnect(new PipeMessageData(pd.Pipe, null, pd.State));
                _connections.TryRemove(pd.Pipe, out IPCData nullobj);
            }
        }

        private void ProcessReceivedMessage(IAsyncResult result)
        {
            // read from client until complete
            IPCData pd = (IPCData)result.AsyncState;
            Int32 bytesRead = pd.Pipe.EndRead(result);
            if (bytesRead > 0)
            {
                pd.DataLength = bytesRead;
                OnMessageReceived(new PipeMessageData(pd.Pipe, pd, pd.State));
            }
            BeginRead(pd);
        }
    }
}
