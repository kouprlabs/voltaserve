namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class GrpcSettings
  {
    [JsonProperty("target")]
    public string Target { get; set; }
  }
}