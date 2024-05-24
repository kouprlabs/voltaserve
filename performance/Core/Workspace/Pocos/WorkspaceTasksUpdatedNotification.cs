namespace Defyle.Core.Workspace.Pocos
{
  using Infrastructure.Poco;

  public class WorkspaceTasksUpdatedNotification : Notification
	{
    public string WorkspaceId { get; set; }
	}
}