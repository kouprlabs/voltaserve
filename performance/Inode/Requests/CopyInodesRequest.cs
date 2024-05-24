namespace Defyle.WebApi.Inode.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class CopyInodesRequest
	{
		[Required]
		public IEnumerable<string> Ids { get; set; }
	}
}