namespace Defyle.WebApi.Inode.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class CreateFileRequest
	{
		[Required]
		public string Name { get; set; }
    
		public string ParentId { get; set; }
  }
}