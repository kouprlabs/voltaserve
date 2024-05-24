namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class ExternalAppSettings
  {
    public class DockerSettings
    {
      [JsonProperty("cli")]
      public string Cli { get; set; }
      
      [JsonProperty("image")]
      public string Image { get; set; }
      
      [JsonProperty("app")]
      public string App { get; set; }

      [JsonProperty("inputDirectory")]
      public string InputDirectory { get; set; }
      
      [JsonProperty("outputDirectory")]
      public string OutputDirectory { get; set; }
    }
    
    [JsonProperty("docker")]
    public DockerSettings Docker { get; set; }

    [JsonProperty("binary")]
    public string App { get; set; }
  }
}