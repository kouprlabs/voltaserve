namespace Defyle.Core.Streaming.Services
{
  using System.IO;
  using System.Threading.Tasks;
  using Infrastructure.Poco;
  using Storage.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class PdfStreamService : StreamService
	{
		private readonly CoreSettings _coreSettings;
    private readonly PathService _pathService;

    public PdfStreamService(CoreSettings coreSettings, PathService pathService)
      : base(coreSettings)
    {
			_coreSettings = coreSettings;
      _pathService = pathService;
    }

    public override async Task<string> GetLocalPathAsync(Workspace workspace, File file)
    {
      return file.GetMime() == "application/pdf" ?
        _pathService.GetOriginalFile(workspace, file) :
        Path.Combine(await GetLocalDirectoryAsync(workspace, file), "document.pdf");
    }

    public override Task<string> GetLocalDirectoryAsync(Workspace workspace, File file) =>
      Task.FromResult(Path.Combine(_coreSettings.DataDirectory, workspace.Id, "file-data", file.Id, "document-preview"));

    public override Task<string> GetS3KeyAsync(Workspace workspace, File file)
    {
      return file.GetMime() == "application/pdf" ?
        Task.FromResult(_pathService.GetS3OriginalFile(workspace, file)) :
        Task.FromResult(Path.Combine(workspace.Id, "file-data", file.Id, "document-preview", "document.pdf"));
    }
  }
}