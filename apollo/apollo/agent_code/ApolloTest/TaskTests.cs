using System;
using System.Collections.Generic;
using Phantom;
using Phantom.Jobs;
using Phantom.Tasks;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace ApolloTests
{
    [TestClass]
    public class TaskTests
    {
    private static void Zc3d4e5()
    {
        Thread.Sleep(Random.Next(1, 5));
        GC.Collect();
    }
        
        
        [TestMethod]
        public void CatTest()
        {
            if (!System.IO.File.Exists("C:\\Users\\Public\\test.txt"))
            {
                using (System.IO.FileStream fs = System.IO.File.Create("C:\\Users\\Public\\test.txt"))
                {
                    using (System.IO.StreamWriter sw = new System.IO.StreamWriter(fs))
                    {
                        sw.WriteLine("test file");
                    }
                }
            }    
            Task task = new Task("cat", "C:\\Users\\Public\\test.txt", "1");
            Job job = new Job(task, null);
            Cat.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual("test file", task.message);
        }
        [TestMethod]
        public void CatTestInvalid()
        {
            Task task = new Task("cat", "C:\\balahsdghaseter.txt", "1");
            Job job = new Job(task, null);
            Cat.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "error");
        }
        [TestMethod]
        public void CdTest()
        {
            
            System.IO.Directory.SetCurrentDirectory("C:\\");
            Task task = new Task("cd", "C:\\Users\\Public", "1");
            Job job = new Job(task, null);
            ChangeDir.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual("C:\\Users\\Public", Environment.CurrentDirectory);
            
            System.IO.Directory.SetCurrentDirectory("C:\\");
        }
        [TestMethod]
        public void CdTestInvalid()
        {
            
            System.IO.Directory.SetCurrentDirectory("C:\\");
            Task task = new Task("cd", "C:\\asdfasdthetherhasdf", "1");
            Job job = new Job(task, null);
            ChangeDir.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "error");
            
            System.IO.Directory.SetCurrentDirectory("C:\\");
        }
        [TestMethod]
        public void CopyTest()
        {
            if (System.IO.File.Exists("C:\\Users\\Public\\test2.txt"))
                System.IO.File.Delete("C:\\Users\\Public\\test2.txt");
            Task task = new Task("cp", "C:\\Users\\Public\\test.txt C:\\Users\\Public\\test2.txt", "1");
            Job job = new Job(task, null);
            Copy.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsTrue(System.IO.File.Exists("C:\\Users\\Public\\test2.txt"));
            System.IO.File.Delete("C:\\Users\\Public\\test2.txt");
        } 
        [TestMethod]
        public void CopyTestInvalid()
        {
            if (System.IO.File.Exists("C:\\Users\\Public\\test3.txt"))
                System.IO.File.Delete("C:\\Users\\Public\\test3.txt");
            Task task = new Task("cp", "C:\\asdfasdfathethiethzscgvnbzxg.aste C:\\Users\\Public\\test3.txt", "1");
            Job job = new Job(task, null);
            Copy.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "error");
        }
        [TestMethod]
        public void DirListTest()
        {
            Task task = new Task("ls", "C:\\", "1");
            Job job = new Job(task, null);
            DirectoryList.K4n5o6p7(job, null);
            Console.WriteLine(task.message.GetType());
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual(true, (task.message is List<Apfell.Structs.FileInformation>));
        }
        [TestMethod]
        public void DirListTestInvalid()
        {
            Task task = new Task("ls", "C:\\ethetrhehtet", "1");
            Job job = new Job(task, null);
            DirectoryList.K4n5o6p7(job, null);
            Console.WriteLine(task.message.GetType());
            
            Assert.IsTrue(task.status == "error");
        }
        
        
        
        [TestMethod]
        public void KillTest()
        {
            int procId = System.Diagnostics.Process.Start("notepad.exe").Id;
            System.Diagnostics.Process proc = System.Diagnostics.Process.GetProcessById(procId);
            Assert.IsTrue(!proc.HasExited);
            Task task = new Task("kill", $"{procId}", "1");
            Job job = new Job(task, null);
            Kill.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsTrue(proc.HasExited);
        }
        [TestMethod]
        public void KillTestInvalid()
        {
            Task task = new Task("kill", "1111111111", "1");
            Job job = new Job(task, null);
            Kill.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "error");
        }
        [TestMethod]
        public void PowerShellTest()
        {
            string command = "Get-Process -Name explorer";
            Task task = new Task("powershell", command, "1");
            Job job = new Job(task, null);
            PowerShellManager.K4n5o6p7(job, new Agent(default));
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual(true, (task.message.ToString().Contains("ProcessName")));
        }
        [TestMethod]
        public void PowerShellTestInvalid()
        {
            string command = "Get-AFDSADSHETHWET";
            Task task = new Task("powershell", command, "1");
            Job job = new Job(task, null);
            PowerShellManager.K4n5o6p7(job, new Agent(default));
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual(true, (task.message.ToString().Contains("ERROR")));
        }
        [TestMethod]
        public void PwdTest()
        {
            Task task = new Task("pwd", null, "1");
            Job job = new Job(task, null);
            PrintWorkingDirectory.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual(true, (task.message.ToString().Contains("C:\\")));
        }
        [TestMethod]
        public void ProcessTest()
        {
            Agent agent = new Agent(default);
            Task task = new Task("run", "whoami /priv", "1");
            Job job = new Job(task, agent);
            Phantom.Tasks.Process.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.AreEqual(true, (task.message.ToString().Contains("Process executed")));
        }
        [TestMethod]
        public void ProcessTestInvalid()
        {
            Agent agent = new Agent(default);
            Task task = new Task("run", "blah /asdf", "1");
            Job job = new Job(task, agent);
            Phantom.Tasks.Process.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "error");
        }
        [TestMethod]
        public void ProcessListTest()
        {
            Task task = new Task("ps", null, "1");
            Job job = new Job(task, null);
            Phantom.Tasks.ProcessList.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsTrue(task.message is List<Apfell.Structs.ProcessEntry>);
        }
        [TestMethod]
        public void RemoveTest()
        {
            System.IO.File.Copy("C:\\Users\\Public\\test.txt", "C:\\Users\\Public\\asdfasdf.txt");
            Task task = new Task("rm", "C:\\Users\\Public\\asdfasdf.txt", "1");
            Job job = new Job(task, null);
            Phantom.Tasks.Remove.K4n5o6p7(job, null);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsFalse(System.IO.File.Exists("C:\\Users\\Public\\asdfasdf.txt"));
        }
        [TestMethod]
        public void StealTokenTest()
        {
            Agent agent = new Agent(default);
            Task task = new Task("steal_token", null, "1");
            Job job = new Job(task, agent);
            Token.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsTrue(agent.HasAlternateToken());
            Token.stolenHandle = IntPtr.Zero;
        } 
        [TestMethod]
        public void StealTokenTestInvalid()
        {
            Agent agent = new Agent(default);
            Task task = new Task("steal_token", "1351251251", "1");
            Job job = new Job(task, agent);
            Token.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "error");
            
            Assert.IsFalse(agent.HasAlternateToken());
        }
        [TestMethod]
        public void RevertTokenTest()
        {
            Agent agent = new Agent(default);
            Task task = new Task("steal_token", null, "1");
            Job job = new Job(task, agent);
            Token.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsTrue(agent.HasAlternateToken());

            task = new Task("rev2self", null, "1");
            job = new Job(task, agent);
            Token.K4n5o6p7(job, agent);
            
            Assert.IsTrue(task.status == "complete");
            
            Assert.IsFalse(agent.HasAlternateToken());
        } 
    }
}
