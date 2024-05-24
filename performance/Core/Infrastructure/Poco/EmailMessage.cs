namespace Defyle.Core.Infrastructure.Poco
{
  using System;

  public class EmailMessage
	{
		public string Id { get; set; }

    public string ToEmail { get; set; }
    
    public string ToName { get; set; }

		public string Subject { get; set; }

    public string PlainTextContent { get; set; }

    public string HtmlContent { get; set; }

    public EmailMessage()
    {
      Id = Guid.NewGuid().ToString();
    }
  }
}