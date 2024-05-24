namespace Defyle.WebApi.Inode.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class UpdateInodeRequest
	{
    [StringLength(4000)]
		public string Name { get; set; }
	}
}