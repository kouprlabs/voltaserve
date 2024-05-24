namespace Defyle.WebApi.Inode.Dtos
{
  using System.Collections.Generic;

  public class InodeFacetDto : InodeDto
	{
    public FileFacetDto File { get; set; }

    public DirectoryFacetDto Directory { get; set; }
    
    public IEnumerable<string> Permissions { get; set; }

    public IEnumerable<InodeDto> Path { get; set; }
	}
}