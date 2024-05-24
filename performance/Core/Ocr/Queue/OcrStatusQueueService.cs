namespace Defyle.Core.Ocr.Queue
{
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Preview.Queue;

  public class OcrStatusQueueService : QueueService
  {
    private readonly OcrWorkerSettings _ocrWorkerSettings;
    
    public OcrStatusQueueService(
      MessageBrokerSettings messageBrokerSettings,
      OcrWorkerSettings ocrWorkerSettings)
      : base(messageBrokerSettings)
    {
      _ocrWorkerSettings = ocrWorkerSettings;
    }
    
    public void SendSearchablePdfStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _ocrWorkerSettings.SearchablePdfStatusQueue);
    }
    
    public void SendTextExtractionStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _ocrWorkerSettings.TextExtractionStatusQueue);
    }
  }
}