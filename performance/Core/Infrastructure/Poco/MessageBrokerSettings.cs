namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class MessageBrokerSettings
	{
    [JsonProperty("url")]
    public string Url { get; set; }
	}
}