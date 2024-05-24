namespace Defyle.WebApi.Inode.Responses
{
  using System.Collections.Generic;
  using Core.Infrastructure.Poco;

  public class CopyInodesResponse
	{
		public List<string> Succeeded { get; set; } = new List<string>();
		public List<string> Failed { get; set; } = new List<string>();
		public List<string> Created { get; set; } = new List<string>();
    public List<Error> Errors { get; set; } = new List<Error>();
	}
}