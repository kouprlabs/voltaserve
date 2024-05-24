namespace Defyle.Core.Inode.Models
{
  public class Inode
  {
    public const string InodeTypeFile = "file";
    public const string InodeTypeDirectory = "directory";
    
    
    public string Id { get; set; }
    
    public string WorkspaceId { get; set; }

    public string Name { get; set; }

    public string Type { get; set; }

    public string ParentId { get; set; }
    
    public long SortOrder { get; set; }
    
    public string Text { get; set; }

    public long CreateTime { get; set; }

    public long UpdateTime { get; set; }
  }
}