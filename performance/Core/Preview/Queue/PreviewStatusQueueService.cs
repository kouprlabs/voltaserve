namespace Defyle.Core.Preview.Queue
{
  using Infrastructure.Poco;
  using Infrastructure.Services;

  public class PreviewStatusQueueService : QueueService
  {
    private readonly PreviewWorkerSettings _previewWorkerSettings;

    public PreviewStatusQueueService(
      MessageBrokerSettings messageBrokerSettings,
      PreviewWorkerSettings previewWorkerSettings)
      : base(messageBrokerSettings)
    {
      _previewWorkerSettings = previewWorkerSettings;
    }

    public void SendImageStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _previewWorkerSettings.ImageStatusQueue);
    }

    public void SendTileMapStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _previewWorkerSettings.TileMapStatusQueue);
    }

    public void SendDocumentStatusMessage(string workspaceId, string inodeId, string previewStatus)
    {
      var message = new PreviewStatusMessage
      {
        WorkspaceId = workspaceId,
        InodeId = inodeId,
        Status = previewStatus
      };
      
      SendQueueMessage(message, _previewWorkerSettings.DocumentStatusQueue);
    }
  }
}