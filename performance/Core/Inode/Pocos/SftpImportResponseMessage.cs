namespace Defyle.Core.Inode.Pocos
{
	public class SftpImportResponseMessage
	{
		public string Id { get; set; }

		public string MessageId { get; set; }

		public string WorkspaceId { get; set; }

		public bool Success { get; set; }

		public string StatusMessage { get; set; }
	}
}