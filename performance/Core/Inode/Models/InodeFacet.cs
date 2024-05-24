namespace Defyle.Core.Inode.Models
{
  using System.Collections.Generic;

  public class InodeFacet : Inode
  {
    public FileFacet File { get; set; }

    public DirectoryFacet Directory { get; set; }
    
    public IEnumerable<string> Permissions { get; set; }

    public IEnumerable<Inode> Path { get; set; }
  }
}