namespace Defyle.Core.Infrastructure.Poco
{
  using System.Collections.Generic;
  using Newtonsoft.Json;

  public class CoreSettings
	{
    [JsonProperty("messageBroker")]
		public MessageBrokerSettings MessageBroker { get; set; }

    [JsonProperty("workspaceInvitationEmail")]
    public EmailSettings WorkspaceInvitationEmail { get; set; }
    
    [JsonProperty("resetPasswordEmail")]
		public EmailSettings ResetPasswordEmail { get; set; }

    [JsonProperty("confirmationEmail")]
		public EmailSettings ConfirmationEmail { get; set; }

    [JsonProperty("password")]
		public PasswordSettings Password { get; set; }

    [JsonProperty("token")]
		public TokenSettings Token { get; set; }

    [JsonProperty("ldap")]
		public LdapSettings Ldap { get; set; }

    [JsonProperty("elasticsearch")]
    public ElasticsearchSettings Elasticsearch { get; set; }
    
    [JsonProperty("coreService")]
    public GrpcSettings CoreService { get; set; }
    
    [JsonProperty("previewWorker")]
    public PreviewWorkerSettings PreviewWorker { get; set; }

    [JsonProperty("ocrWorker")]
    public OcrWorkerSettings OcrWorker { get; set; }

    [JsonProperty("webClientUrl")]
    public string WebClientUrl { get; set; }

    [JsonProperty("authenticationType")]
		public string AuthenticationType { get; set; }

    [JsonProperty("partitionId")]
		public string PartitionId { get; set; }

    [JsonProperty("url")]
		public string Url { get; set; }

    private string _dataDirectory;
      
    [JsonProperty("dataDirectory")]
    public string DataDirectory
    {
      get => _dataDirectory;
      set => _dataDirectory = PathUtils.Rewrite(value);
    }
    
    private string _tempDirectory;

    [JsonProperty("tempDirectory")]
    public string TempDirectory
    {
      get => _tempDirectory;
      set => _tempDirectory = PathUtils.Rewrite(value);
    }

    [JsonProperty("allowedCorsOrigins")]
    public List<string> AllowedCorsOrigins { get; set; } = new List<string>();

    [JsonProperty("allowedHosts")]
		public List<string> AllowedHosts { get; set; } = new List<string>();

    [JsonProperty("cleanupQueue")]
    public string CleanupQueue { get; set; }
    
    [JsonProperty("emailQueue")]
		public string EmailQueue { get; set; }

    [JsonProperty("sftpQueue")]
		public string SftpQueue { get; set; }

    [JsonProperty("s3Queue")]
    public string S3Queue { get; set; }
    
    [JsonProperty("s3")]
    public S3Settings S3 { get; set; }
    
    [JsonProperty("s3Enabled")]
    public bool S3Enabled { get; set; }

    [JsonProperty("enableTileMapPreview")]
		public bool EnableTileMapPreview { get; set; }

    [JsonProperty("multipartBodyLengthLimit")]
		public long MultipartBodyLengthLimit { get; set; }
    
    [JsonProperty("securityKey")]
    public string SecurityKey { get; set; }

    [JsonProperty("gravatarIntegration")]
		public bool GravatarIntegration { get; set; }
    
    [JsonProperty("defaultWorkspaceStorageCapacity")]
		public long DefaultWorkspaceStorageCapacity { get; set; }

    [JsonProperty("previewFileSizeLimit")]
    public long PreviewFileSizeLimit { get; set; }
  }
}