namespace Defyle.Core.Streaming.Services
{
  using System.IO;
  using System.Threading.Tasks;
  using Infrastructure.Poco;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class OcrPdfStreamService : StreamService
	{
		private readonly CoreSettings _coreSettings;

    public OcrPdfStreamService(CoreSettings coreSettings)
      : base(coreSettings)
    {
			_coreSettings = coreSettings;
    }

    public override async Task<string> GetLocalPathAsync(Workspace workspace, File file)
    {
      return Path.Combine(await GetLocalDirectoryAsync(workspace, file), "document.pdf");
    }

    public override Task<string> GetLocalDirectoryAsync(Workspace workspace, File file) =>
      Task.FromResult(Path.Combine(_coreSettings.DataDirectory, workspace.Id, "file-data", file.Id, "ocr"));

    public override Task<string> GetS3KeyAsync(Workspace workspace, File file) =>
      Task.FromResult(Path.Combine(workspace.Id, "file-data", file.Id, "ocr", "document.pdf"));
  }
}