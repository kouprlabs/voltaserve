namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class EmailSettings
	{
    [JsonProperty("subject")]
		public string Subject { get; set; }

    [JsonProperty("htmlContent")]
		public string[] HtmlContent { get; set; }
    
    [JsonProperty("plainTextContent")]
    public string[] PlainTextContent { get; set; }
    
    public string GetHtmlContent() => string.Join("", HtmlContent);
    
    public string GetPlainTextContent() => string.Join("\n", PlainTextContent);
  }
}