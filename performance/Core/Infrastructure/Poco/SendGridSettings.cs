namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class SendGridSettings
  {
    [JsonProperty("apiKey")]
    public string ApiKey { get; set; }
  }
}