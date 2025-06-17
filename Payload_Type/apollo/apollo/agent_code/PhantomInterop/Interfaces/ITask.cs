using System.Threading.Tasks;
namespace PhantomInterop.Interfaces
{
    public interface ITask
    {
        string ID();
        void Start();
        Task CreateTasking();
        void Kill();
    }
}
