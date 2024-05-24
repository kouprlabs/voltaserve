namespace Defyle.Core.Inode.Services
{
  using System;
  using System.Collections.Generic;
  using System.Threading.Tasks;
  using Infrastructure.Services;
  using Microsoft.AspNetCore.SignalR;
  using Newtonsoft.Json;
  using Pocos;

  public class InodeNotificationService
  {
    private readonly IHubContext<NotificationHub> _hubContext;
    private readonly JsonService _jsonService;

    public InodeNotificationService(
      IHubContext<NotificationHub> hubContext,
      JsonService jsonService)
    {
      _hubContext = hubContext;
      _jsonService = jsonService;
    }
    
    public async Task SendInodesChildrenUpdatedAsync(string workspaceId, IEnumerable<string> ids)
    {
      var notification = new InodesChildrenUpdatedNotification
      {
        MessageId = Guid.NewGuid().ToString(),
        WorkspaceId = workspaceId,
        Ids = ids
      };

      await _hubContext.Clients.All.SendAsync(
        "InodesChildrenUpdated",
        JsonConvert.SerializeObject(notification, _jsonService.GetJsonSerializerSettings()));
    }
    
    public async Task SendInodesPropertiesUpdatedAsync(string workspaceId, IEnumerable<string> ids)
    {
      var notification = new InodesPropertiesUpdatedNotification
      {
        MessageId = Guid.NewGuid().ToString(),
        WorkspaceId = workspaceId,
        Ids = ids
      };

      await _hubContext.Clients.All.SendAsync(
        "InodesPropertiesUpdated",
        JsonConvert.SerializeObject(notification, _jsonService.GetJsonSerializerSettings()));
    }
  }
}