namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class ElasticsearchSettings
  {
    [JsonProperty("url")]
    public string Url { get; set; }

    [JsonProperty("inodesIndex")]
    public string InodesIndex { get; set; }
  }
}