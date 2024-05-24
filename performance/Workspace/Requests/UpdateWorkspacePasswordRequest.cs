namespace Defyle.WebApi.Workspace.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class UpdateWorkspacePasswordRequest
	{
		[Required]
		public string CurrentPassword { get; set; }

		[Required]
		public string NewPassword { get; set; }
	}
}