namespace Defyle.Core.Workspace.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class WorkspaceNotification : Notification
	{
    public string Id { get; set; }
    
		public string WorkspaceId { get; set; }
    
		public string Title { get; set; }
    
		public IEnumerable<string> Content { get; set; }
    
    public long CreatedAt { get; set; }
	}
}