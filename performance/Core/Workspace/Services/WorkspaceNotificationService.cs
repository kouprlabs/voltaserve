namespace Defyle.Core.Workspace.Services
{
  using System;
  using System.Collections.Generic;
  using System.Threading.Tasks;
  using Infrastructure.Services;
  using Microsoft.AspNetCore.SignalR;
  using Newtonsoft.Json;
  using Pocos;

  public class WorkspaceNotificationService
  {
    private readonly IHubContext<NotificationHub> _hubContext;

    public WorkspaceNotificationService(IHubContext<NotificationHub> hubContext)
    {
      _hubContext = hubContext;
    }

    public async Task SendWorkspaceNotificationAsync(string workspaceId, string title, IEnumerable<string> content)
    {
      var notification = new WorkspaceNotification
      {
        MessageId = Guid.NewGuid().ToString(),
        Title = title,
        Content = content,
        WorkspaceId = workspaceId,
        CreatedAt = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()
      };
      await _hubContext.Clients.All.SendAsync("WorkspaceNotification", JsonConvert.SerializeObject(notification));
    }

    public async Task SendWorkspacesUpdatedAsync()
    {
      var notification = new WorkspacesUpdatedNotification
      {
        MessageId = Guid.NewGuid().ToString()
      };
      await _hubContext.Clients.All.SendAsync("WorkspacesUpdated", JsonConvert.SerializeObject(notification));
    }

    public async Task SendWorkspaceTasksUpdatedAsync(string workspaceId)
    {
      var notification = new WorkspaceTasksUpdatedNotification
      {
        MessageId = Guid.NewGuid().ToString(),
        WorkspaceId = workspaceId
      };
      
      await _hubContext.Clients.All.SendAsync("WorkspaceTasksUpdated", JsonConvert.SerializeObject(notification));
    }
  }
}