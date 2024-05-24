namespace Defyle.Core.Storage.Models
{
  using System.Collections.Generic;

  public class FileCategory
  {
    public string Id { get; set; }

    public List<string> PhysicalFileTypes { get; set; }
  }
}