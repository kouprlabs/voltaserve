namespace Defyle.Core.Inode.Models
{
  using System.Collections.Generic;
  using Storage.Models;

  public class FileFacet
  {
    public string Mime { get; set; }
    
    public long Size { get; set; }

    public long Version { get; set; }

    public string Category { get; set; }

    public string Type { get; set; }

    public IEnumerable<FileProperty> Properties { get; set; }
  }
}