namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class SmtpServerSettings
	{
    [JsonProperty("host")]
		public string Host { get; set; }
    
    [JsonProperty("port")]
		public int Port { get; set; }
    
    [JsonProperty("username")]
		public string Username { get; set; }
    
    [JsonProperty("password")]
		public string Password { get; set; }

    [JsonProperty("useDefaultCredentials")]
    public bool? UseDefaultCredentials { get; set; }

    [JsonProperty("deliveryMethod")]
    public string DeliveryMethod { get; set; }

    [JsonProperty("deliveryFormat")]
    public string DeliveryFormat { get; set; }

    [JsonProperty("timeout")]
    public int? Timeout { get; set; }
    
    [JsonProperty("enableSsl")]
    public bool? EnableSsl { get; set; }

    [JsonProperty("targetName")]
    public string TargetName { get; set; }
  }
}