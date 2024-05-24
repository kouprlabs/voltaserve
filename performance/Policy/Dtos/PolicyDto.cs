namespace Defyle.WebApi.Policy.Dtos
{
  public class PolicyDto
  {
    public string Id { get; set; }
    
    public string Subject { get; set; }

    public string Object { get; set; }

    public string Permission { get; set; }
    
    public long CreateTime { get; set; }
  }
}