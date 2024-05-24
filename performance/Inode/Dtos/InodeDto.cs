namespace Defyle.WebApi.Inode.Dtos
{
  public class InodeDto
  {
    public string Id { get; set; }
    
    public string WorkspaceId { get; set; }

    public string Name { get; set; }

    public string Type { get; set; }

    public string ParentId { get; set; }
    
    public long SortOrder { get; set; }

    public long CreateTime { get; set; }

    public long UpdateTime { get; set; }
  }
}