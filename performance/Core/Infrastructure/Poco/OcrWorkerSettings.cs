namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class OcrWorkerSettings
	{
    [JsonProperty("textExtractionWorkQueue")]
    public string TextExtractionWorkQueue { get; set; }

    [JsonProperty("textExtractionStatusQueue")]
    public string TextExtractionStatusQueue { get; set; }
    
    [JsonProperty("searchablePdfWorkQueue")]
    public string SearchablePdfWorkQueue { get; set; }

    [JsonProperty("searchablePdfStatusQueue")]
    public string SearchablePdfStatusQueue { get; set; }
	}
}