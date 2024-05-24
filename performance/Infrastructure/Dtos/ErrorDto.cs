namespace Defyle.WebApi.Infrastructure.Dtos
{
  using Newtonsoft.Json;

  public class ErrorDto
  {
    [JsonProperty("code")]
    public string Code { get; set; }

    [JsonProperty("description")]
    public string Description { get; set; }
  }
}