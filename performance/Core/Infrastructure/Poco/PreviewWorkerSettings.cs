namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class PreviewWorkerSettings
	{
    [JsonProperty("imageWorkQueue")]
		public string ImageWorkQueue { get; set; }

    [JsonProperty("imageStatusQueue")]
		public string ImageStatusQueue { get; set; }

    [JsonProperty("tileMapWorkQueue")]
		public string TileMapWorkQueue { get; set; }

    [JsonProperty("tileMapStatusQueue")]
		public string TileMapStatusQueue { get; set; }

    [JsonProperty("documentWorkQueue")]
		public string DocumentWorkQueue { get; set; }

    [JsonProperty("documentStatusQueue")]
		public string DocumentStatusQueue { get; set; }
  }
}