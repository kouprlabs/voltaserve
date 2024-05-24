namespace Defyle.Core.Streaming.Services
{
  using System.Threading.Tasks;
  using Infrastructure.Poco;
  using Storage.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class OriginalStreamService : StreamService
  {
    private readonly PathService _pathService;

    public OriginalStreamService(CoreSettings coreSettings, PathService pathService)
      : base(coreSettings)
    {
      _pathService = pathService;
    }

    public override Task<string> GetLocalPathAsync(Workspace workspace, File file) =>
      Task.FromResult(_pathService.GetOriginalFile(workspace, file));

    public override Task<string> GetLocalDirectoryAsync(Workspace workspace, File file) =>
      Task.FromResult(_pathService.GetFileDirectory(workspace, file));

    public override Task<string> GetS3KeyAsync(Workspace workspace, File file) =>
      Task.FromResult(_pathService.GetS3OriginalFile(workspace, file));
  }
}