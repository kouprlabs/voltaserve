namespace Defyle.Core.Workspace.Models
{
  public class WorkspaceTask
  {
    public const string WorkspaceTaskTypeIndeterminateProgress = "indeterminate-progress";
    
		public string Id { get; set; }

		public string Type { get; set; }

		public string Title { get; set; }

		public string Content { get; set; }

    public long CreatedAt { get; set; }
	}
}