namespace Defyle.WebApi.Inode.Responses
{
  using System.Collections.Generic;
  using Core.Infrastructure.Poco;

  public class DeleteInodesResponse
	{
		public List<string> AffectedParentIds { get; set; } = new List<string>();
		public List<string> Succeeded { get; set; } = new List<string>();
		public List<string> Failed { get; set; } = new List<string>();
    public List<Error> Errors { get; set; } = new List<Error>();
	}
}