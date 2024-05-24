namespace Defyle.WebApi.Infrastructure.Responses
{
  using System.Collections.Generic;
  using Dtos;
  using Newtonsoft.Json;

  public class ErrorsResponse
  {
    [JsonProperty("errors")]
    public IEnumerable<ErrorDto> Errors { get; set; }
  }
}