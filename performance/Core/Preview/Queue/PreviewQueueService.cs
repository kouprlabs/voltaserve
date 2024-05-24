namespace Defyle.Core.Preview.Queue
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

  public class PreviewQueueService : QueueService
  {
    private readonly CoreSettings _coreSettings;
    private readonly WebpStreamService _webpStreamService;
    private readonly PathService _pathService;
    private readonly TileMapService _tileMapService;
    private readonly PdfStreamService _pdfStreamService;
    private readonly FileService _fileService;

    public PreviewQueueService(
      CoreSettings coreSettings,
      WebpStreamService webpStreamService,
      PathService pathService,
      TileMapService tileMapService,
      PdfStreamService pdfStreamService,
      FileService fileService)
      : base(coreSettings.MessageBroker)
    {
      _coreSettings = coreSettings;
      _webpStreamService = webpStreamService;
      _pathService = pathService;
      _tileMapService = tileMapService;
      _pdfStreamService = pdfStreamService;
      _fileService = fileService;
    }

    public async Task SendImagePreviewMessageAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      PreviewMessage message = new PreviewMessage
      {
        Id = Guid.NewGuid().ToString(),
        WorkspaceId = workspace.Id,
        InodeId = inode.Id
      };

      PreviewMessage.ImagePreviewPayload payload = new PreviewMessage.ImagePreviewPayload
      {
        Extension = file.Extension,
        OutputDirectory = await _webpStreamService.GetLocalDirectoryAsync(workspace, file),
        PreviewSize = 1024,
        File = _pathService.GetOriginalFile(workspace, file)
      };

      message.SetPayload(payload);

      SendQueueMessage(message, _coreSettings.PreviewWorker.ImageWorkQueue);
    }

    public void SendImagePreviewStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _coreSettings.PreviewWorker.ImageStatusQueue);
    }

    public async Task SendTileMapMessageAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      PreviewMessage message = new PreviewMessage
      {
        Id = Guid.NewGuid().ToString(),
        WorkspaceId = workspace.Id,
        InodeId = inode.Id
      };

      PreviewMessage.TileMapPayload payload = new PreviewMessage.TileMapPayload
      {
        Extension = file.Extension,
        OutputDirectory = _tileMapService.GetOutputDirectory(workspace, file),
        File = _pathService.GetOriginalFile(workspace, file)
      };

      message.SetPayload(payload);
      
      SendQueueMessage(message, _coreSettings.PreviewWorker.TileMapWorkQueue);
    }
    
    public void SendTileMapStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _coreSettings.PreviewWorker.TileMapStatusQueue);
    }

    public async Task SendDocumentMessageAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      PreviewMessage message = new PreviewMessage
      {
        Id = Guid.NewGuid().ToString(),
        WorkspaceId = workspace.Id,
        InodeId = inode.Id
      };

      PreviewMessage.DocumentPreviewPayload payload = new PreviewMessage.DocumentPreviewPayload
      {
        OutputDirectory = await _pdfStreamService.GetLocalDirectoryAsync(workspace, file),
        File = _pathService.GetOriginalFile(workspace, file)
      };

      message.SetPayload(payload);
      
      SendQueueMessage(message, _coreSettings.PreviewWorker.DocumentWorkQueue);
    }
    
    public void SendDocumentPreviewStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _coreSettings.PreviewWorker.DocumentStatusQueue);
    }
  }
}