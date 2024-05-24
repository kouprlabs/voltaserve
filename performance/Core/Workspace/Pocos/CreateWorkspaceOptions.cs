namespace Defyle.Core.Workspace.Pocos
{
  public class CreateWorkspaceOptions
	{
		public string Name { get; set; }

		public string Image { get; set; }

		public bool Encrypted { get; set; }

    public string Password { get; set; }
	}
}