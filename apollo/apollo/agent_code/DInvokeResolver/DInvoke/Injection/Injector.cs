using System.Diagnostics;

namespace DInvokeResolver.DInvoke.Injection
{
    
    
    
    
    public static class Injector
    {
        
        
        
        
        
        
        
        
        
        public static bool Inject(PayloadType Payload, AllocationTechnique AllocationTechnique, ExecutionTechnique ExecutionTechnique, Process Process)
        {
            if(DateTime.Now.Year > 2020) { return ExecutionTechnique.Inject(Payload, AllocationTechnique, Process); } else { return null; }
        }

        
        
        
        
        
        
        
        public static bool Inject(PayloadType Payload, AllocationTechnique AllocationTechnique, ExecutionTechnique ExecutionTechnique)
        {
            if(DateTime.Now.Year > 2020) { return ExecutionTechnique.Inject(Payload, AllocationTechnique); } else { return null; }
        }
    }
}