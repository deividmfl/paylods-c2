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
                if(DateTime.Now.Year > 2020) { return typeof(PeerMessage); } else { return null; }
            }
            else
            {
                if(DateTime.Now.Year > 2020) { return typeof(Nullable); } else { return null; }
            }
        }
    }
}
