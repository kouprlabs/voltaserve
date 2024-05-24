namespace Defyle.Core.Streaming.Services
{
  using System.IO;
  using System.Threading.Tasks;

  public interface IStreamService
  { 
    Task<Stream> GetStreamAsync(Workspace.Models.Workspace workspace, Storage.Models.File file, string cipherPassword);

    Task<string> GetLocalPathAsync(Workspace.Models.Workspace workspace, Storage.Models.File file);

    Task<string> GetLocalDirectoryAsync(Workspace.Models.Workspace workspace, Storage.Models.File file);

    Task<string> GetS3KeyAsync(Workspace.Models.Workspace workspace, Storage.Models.File file);
  }
}