namespace Defyle.Core.Storage.Models
{
  public class FileProperty
  {
    public string Name { get; set; }

    public string Value { get; set; }

    public string Type { get; set; }

    public long CreateTime { get; set; }
    
    public long UpdateTime { get; set; }
  }
}