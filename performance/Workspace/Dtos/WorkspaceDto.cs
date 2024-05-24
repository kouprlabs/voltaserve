namespace Defyle.WebApi.Workspace.Dtos
{
  using System.Collections.Generic;

  public class WorkspaceDto
	{
		public string Id { get; set; }

		public string PartitionId { get; set; }

		public string Name { get; set; }

		public string Image { get; set; }

    public IEnumerable<string> Permissions { get; set; }

    public bool Encrypted { get; set; }

    public IEnumerable<WorkspaceTaskDto> Tasks { get; set; }

    public long StorageCapacity { get; set; }

		public long CreatedTime { get; set; }

		public long UpdatedTime { get; set; }
	}
}