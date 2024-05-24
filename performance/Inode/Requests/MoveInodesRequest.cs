namespace Defyle.WebApi.Inode.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class MoveInodesRequest
	{
		[Required]
		public IEnumerable<string> Ids { get; set; }
	}
}