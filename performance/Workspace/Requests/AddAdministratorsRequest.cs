namespace Defyle.WebApi.Workspace.Requests
{
  using System.Collections.Generic;

  public class AddAdministratorsRequest
	{
		public IEnumerable<string> Emails { get; set; }
	}
}