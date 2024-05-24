namespace Defyle.Core.Infrastructure.Poco
{
	public class Error
	{
    public static readonly Error MissingPartitionIdHeaderError = new Error(
			"MissingPartitionIdHeader",
			"Missing partition id header.");

		public static readonly Error InvalidEmailError = new Error(
			"InvalidEmail",
			"Invalid email.");
    
    public static readonly Error InvalidFullNameError = new Error(
      "InvalidFullName",
      "Invalid full name.");

    public static readonly Error DirectoryOnlyOperationError = new Error(
			"DirectoryOnlyOperation",
			"This operation is supported on directories only.");

    public static readonly Error InodeNotFoundError = new Error(
			"InodeNotFound",
			"Inode not found.");

    public static readonly Error PreviewNotFoundError = new Error(
      "PreviewNotFound",
      "Preview not found.");

    public static readonly Error PhysicalFileNotFoundError = new Error(
			"PhysicalFileNotFound",
			"Physical file not found.");

    public static readonly Error NullOrEmptyPasswordError = new Error(
      "NullOrEmptyPasswordError",
      "Password is null or empty.");
    
    public static readonly Error InvalidEmailOrPasswordError = new Error(
			"InvalidEmailOrPassword",
			"Invalid email or password.");

		public static readonly Error InvalidWorkspacePasswordError = new Error(
			"InvalidWorkspacePassword",
			"Invalid workspace password.");

		public static readonly Error EmailNotConfirmedError = new Error(
			"EmailNotConfirmedError",
			"Email not confirmed.");

		public static readonly Error StorageUsageExceededError = new Error(
			"StorageUsageExceeded",
			"Storage usage exceeded.");

		public static readonly Error LdapConnectionError = new Error(
			"LdapConnectionError",
			"LDAP connection error.");

		public static readonly Error LdapBindError = new Error(
			"LdapBindError",
			"LDAP bind error.");

		public static readonly Error LdapAuthenticationError = new Error(
			"LdapAuthenticationError",
			"LDAP authentication error.");

    public static readonly Error WorkspaceNotEncryptedError = new Error(
			"WorkspaceNotEncrypted",
			"Workspace not encrypted.");

		public static readonly Error WorkspacePasswordRequiredError = new Error(
			"WorkspacePasswordRequired",
			"Workspace password required.");

    public Error()
		{
		}

		public Error(string code, string description)
		{
			Code = code;
			Description = description;
		}

		public string Code { get; }

		public string Description { get; }

		public Error WithCode(string code) => new Error(code, Description);

		public Error WithDescription(string description) => new Error(Code, description);
	}
}