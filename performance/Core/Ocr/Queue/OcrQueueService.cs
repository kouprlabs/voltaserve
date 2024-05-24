namespace Defyle.Core.Ocr.Queue
{
  using System;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Inode.Models;
  using Services;
  using Storage.Models;
  using Storage.Services;
  using Streaming.Services;
  using Workspace.Models;

  public class OcrQueueService : QueueService
  {
    private readonly CoreSettings _coreSettings;
    private readonly PathService _pathService;
    private readonly PdfStreamService _pdfStreamService;
    private readonly OcrPdfStreamService _ocrPdfStreamService;
    private readonly TextExtractionService _textExtractionService;
    private readonly FileService _fileService;

    public OcrQueueService(
      CoreSettings coreSettings,
      PathService pathService,
      PdfStreamService pdfStreamService,
      OcrPdfStreamService ocrPdfStreamService,
      TextExtractionService textExtractionService,
      FileService fileService)
      : base(coreSettings.MessageBroker)
    {
      _coreSettings = coreSettings;
      _pathService = pathService;
      _pdfStreamService = pdfStreamService;
      _ocrPdfStreamService = ocrPdfStreamService;
      _textExtractionService = textExtractionService;
      _fileService = fileService;
    }
    
    public async Task SendSearchablePdfMessageAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      SearchablePdfMessage message = new SearchablePdfMessage
      {
        Id = Guid.NewGuid().ToString(),
        WorkspaceId = workspace.Id,
        InodeId = inode.Id
      };

      string fileType = await _fileService.GetFileTypeAsync(file.Mime);
      string fileCategory = await _fileService.GetFileCategoryAsync(fileType);

      string path;
      if (fileCategory == "document")
      {
        path = await _pdfStreamService.GetLocalPathAsync(workspace, file);
      }
      else if (fileCategory == "image")
      {
        path = _pathService.GetOriginalFile(workspace, file);
      }
      else
      {
        throw new Exception("Unsupported file category.");
      }

      TextExtractionMessage.TextExtractionPayload payload = new TextExtractionMessage.TextExtractionPayload
      {
        OutputDirectory = await _ocrPdfStreamService.GetLocalDirectoryAsync(workspace, file),
        File = path
      };

      message.SetPayload(payload);
      
      SendQueueMessage(message, _coreSettings.OcrWorker.SearchablePdfWorkQueue);
    }

    public async Task SendTextExtractionMessageAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      TextExtractionMessage message = new TextExtractionMessage
      {
        Id = Guid.NewGuid().ToString(),
        WorkspaceId = workspace.Id,
        InodeId = inode.Id
      };
      
      string fileType = await _fileService.GetFileTypeAsync(file.Mime);
      string fileCategory = await _fileService.GetFileCategoryAsync(fileType);

      string path;
      if (file.IndexContent && (fileCategory == "image" || file.GetMime() == "application/pdf"))
      {
        path = await _ocrPdfStreamService.GetLocalPathAsync(workspace, file);
      }
      else
      {
        path = await _pdfStreamService.GetLocalPathAsync(workspace, file);
      }

      TextExtractionMessage.TextExtractionPayload payload = new TextExtractionMessage.TextExtractionPayload
      {
        OutputDirectory = _textExtractionService.GetOutputDirectory(workspace, file),
        File = path
      };

      message.SetPayload(payload);
      
      SendQueueMessage(message, _coreSettings.OcrWorker.TextExtractionWorkQueue);
    }
  }
}