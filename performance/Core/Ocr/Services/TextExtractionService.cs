namespace Defyle.Core.Ocr.Services
{
  using System.IO;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Inode.Models;
  using Storage.Services;
  using Streaming.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class TextExtractionService
	{
		private readonly CoreSettings _coreSettings;
    private readonly PathService _pathService;
    private readonly PdfStreamService _pdfStreamService;
    private readonly OcrPdfStreamService _ocrPdfStreamService;
    private readonly FileService _fileService;

    public TextExtractionService(
			CoreSettings coreSettings,
			PathService pathService,
      PdfStreamService pdfStreamService,
      OcrPdfStreamService ocrPdfStreamService,
      FileService fileService)
    {
			_coreSettings = coreSettings;
      _pathService = pathService;
      _pdfStreamService = pdfStreamService;
      _ocrPdfStreamService = ocrPdfStreamService;
      _fileService = fileService;
    }
    
    public async Task<(Stream stream, string mime, string name)> GetTextFileAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
      }

      string path = await GetPathAsync(workspace, inode, user);

      if (!System.IO.File.Exists(path))
      {
        throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
      }

      if (workspace.Encrypted)
      {
        throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
      }
			
      return (new FileStream(path, FileMode.Open), "text/plain", inode.Name);
    }

    public async Task<string> GetPathAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return null;
      }
      
      return file.GetMime() == "text/plain" ? _pathService.GetOriginalFile(workspace, file) : GetOutputFile(workspace, file);
    }
    
    public string GetOutputDirectory(Workspace workspace, File file)
    {
      return Path.Combine(
        _coreSettings.DataDirectory,
        workspace.Id,
        "file-data",
        file.Id,
        "text");
    }

		private string GetOutputFile(Workspace workspace, File file)
		{
			return Path.Combine(
				GetOutputDirectory(workspace, file),
				"document.txt");
		}
	}
}