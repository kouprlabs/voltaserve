namespace Defyle.Core.Policy.Models
{
  public class Policy
  {
    public string Id { get; set; }
    
    public string Subject { get; set; }

    public string Object { get; set; }

    public string Permission { get; set; }
    
    public long CreateTime { get; set; }
  }
}