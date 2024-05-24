namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class JsonError
  {
    [JsonProperty("code")]
    public string Code { get; set; }

    [JsonProperty("internalDescription")]
    public string InternalDescription { get; set; }
    
    [JsonProperty("publicDescription")]
    public string PublicDescription { get; set; }
  }
}