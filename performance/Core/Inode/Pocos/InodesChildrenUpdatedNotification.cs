namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class InodesChildrenUpdatedNotification : Notification
	{
    public string WorkspaceId { get; set; }
    
		public IEnumerable<string> Ids { get; set; }
	}
}