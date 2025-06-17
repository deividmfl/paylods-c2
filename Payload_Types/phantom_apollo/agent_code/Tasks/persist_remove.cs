using System;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using Phantom.Management.Persistence;

namespace Tasks
{
    public class persist_remove : ICommand
    {
        public string command { get; } = "persist_remove";
        public string description { get; } = "Remove all persistence mechanisms";
        public string help { get; } = "persist_remove";
        public string version { get; } = "1.0.0.0";
        public string id { get; } = "persist_remove";

        public PhantomTaskResponse CreateTasking(PhantomTaskMessage task)
        {
            try
            {
                return PersistenceHandler.RemovePersistence();
            }
            catch (Exception ex)
            {
                return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error removing persistence: {ex.Message}"
                };
            }
        }
    }
}