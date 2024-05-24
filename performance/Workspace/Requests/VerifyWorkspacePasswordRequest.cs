namespace Defyle.WebApi.Workspace.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class VerifyWorkspacePasswordRequest
	{
		[Required]
		public string Password { get; set; }
	}
}