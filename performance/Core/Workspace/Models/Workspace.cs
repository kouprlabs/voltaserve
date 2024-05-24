namespace Defyle.Core.Workspace.Models
{
  using System.Collections.Generic;

  public class Workspace
	{    
		public string Id { get; set; }

		public string PartitionId { get; set; }

    public string RootInodeId { get; set; }

		public string Name { get; set; }

		public string Image { get; set; }

    public bool Encrypted { get; set; }

    public string PasswordHash { get; set; }

		public string CipherKey { get; set; }

		public string TransitKey { get; set; }

		public string TransitIv { get; set; }

		public string Salt { get; set; }

		public List<WorkspaceTask> Tasks { get; set; } = new List<WorkspaceTask>();
    
    public long StorageCapacity { get; set; }

    public long CreateTime { get; set; }
    
		public long UpdateTime { get; set; }
  }
}