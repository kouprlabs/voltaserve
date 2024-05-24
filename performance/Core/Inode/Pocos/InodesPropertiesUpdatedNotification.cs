namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class InodesPropertiesUpdatedNotification : Notification
	{
    public string WorkspaceId { get; set; }
    
		public IEnumerable<string> Ids { get; set; }
	}
}