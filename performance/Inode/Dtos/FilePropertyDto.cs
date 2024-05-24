namespace Defyle.WebApi.Inode.Dtos
{
  public class FilePropertyDto
  {
    public string Name { get; set; }

    public string Value { get; set; }

    public string Type { get; set; }

    public long CreateTime { get; set; }
    
    public long UpdateTime { get; set; }
  }
}