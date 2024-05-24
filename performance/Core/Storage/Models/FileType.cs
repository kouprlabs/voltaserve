namespace Defyle.Core.Storage.Models
{
  using System.Collections.Generic;

  public class FileType
  {
    public string Id { get; set; }

    public List<string> MimeTypes { get; set; }
  }
}