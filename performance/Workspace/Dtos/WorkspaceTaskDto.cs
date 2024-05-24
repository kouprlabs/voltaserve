namespace Defyle.WebApi.Workspace.Dtos
{
	public class WorkspaceTaskDto
	{
		public string Id { get; set; }

		public string Type { get; set; }

		public string Title { get; set; }

		public string Content { get; set; }

    public long CreatedAt { get; set; }
	}
}