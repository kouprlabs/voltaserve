namespace Defyle.WebApi.Inode.Dtos
{
  using System.Collections.Generic;

  public class FileFacetDto
  {
    public string Mime { get; set; }
    
    public long Size { get; set; }
    
    public long Version { get; set; }

    public string Category { get; set; }

    public string Type { get; set; }

    public IEnumerable<FilePropertyDto> Properties { get; set; }
  }
}