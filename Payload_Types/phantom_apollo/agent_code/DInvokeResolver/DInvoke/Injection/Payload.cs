using System;

namespace DInvokeResolver.DInvoke.Injection
{
    
    
    
    
    
    public abstract class PayloadType
    {
        public byte[] Payload { get; private set; }

        
        protected PayloadType(byte[] data)
        {
            Payload = data;
        }
    }

    
    
    
    
    public class PICPayload : PayloadType
    {
        
        public PICPayload(byte[] data) : base(data) { }
    }

    
    
    
    
    public class PayloadTypeNotSupported : Exception
    {
        public PayloadTypeNotSupported() { }

        public PayloadTypeNotSupported(Type payloadType) : base(string.Format("Unsupported Payload type: {0}", payloadType.Name)) { }
    }
}