namespace Defyle.Core.Inode.Pocos
{
  public class SftpImportOptions
	{
    public string Host { get; set; }

		public string Username { get; set; }

		public string Password { get; set; }

		public int Port { get; set; }

		public string Directory { get; set; }

    public bool IndexContent { get; set; }
	}
}