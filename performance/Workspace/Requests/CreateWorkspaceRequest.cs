namespace Defyle.WebApi.Workspace.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class CreateWorkspaceRequest
	{
		[Required]
		[StringLength(100)]
		public string Name { get; set; }

		public string Image { get; set; }

		public bool Encrypted { get; set; }

    public string Password { get; set; }
	}
}