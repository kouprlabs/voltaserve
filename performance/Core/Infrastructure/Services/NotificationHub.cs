namespace Defyle.Core.Infrastructure.Services
{
  using System.Threading.Tasks;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.SignalR;

  [Authorize]
	public class NotificationHub : Hub
	{
		public async Task InodesPropertiesUpdated(string data)
		{
			await Clients.All.SendAsync(nameof(InodesPropertiesUpdated), data);
		}

		public async Task InodesChildrenUpdated(string data)
		{
			await Clients.All.SendAsync(nameof(InodesChildrenUpdated), data);
		}

    public async Task WorkspacesUpdated()
		{
			await Clients.All.SendAsync(nameof(WorkspacesUpdated));
		}

		public async Task WorkspaceTasksUpdated(string data)
		{
			await Clients.All.SendAsync(nameof(WorkspaceTasksUpdated), data);
		}

		public async Task WorkspaceNotification(string data)
		{
			await Clients.All.SendAsync(nameof(WorkspaceNotification), data);
		}
	}
}