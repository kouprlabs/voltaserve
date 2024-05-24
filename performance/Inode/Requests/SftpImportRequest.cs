namespace Defyle.WebApi.Inode.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class SftpImportRequest
	{
		[Required]
		public string Host { get; set; }

		public string Username { get; set; }

		public string Password { get; set; }

		public int Port { get; set; }

		public string Directory { get; set; }

    public bool IndexContent { get; set; }
	}
}