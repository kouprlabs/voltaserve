namespace Defyle.Core.Inode.Pocos
{
  public class SftpImportMessage
	{
		public string Id { get; set; }

		public SftpImportOptions Request { get; set; }

		public string TempDirectory { get; set; }

		public string UserId { get; set; }

		public string WorkspaceId { get; set; }

		public bool IsWorkspaceEncrypted { get; set; }

		public string ParentNodeId { get; set; }

		public bool InheritAccessRights { get; set; }

    public string CipherPassword { get; set; }

		public string ApiUrl { get; set; }

		public string PartitionId { get; set; }
	}
}