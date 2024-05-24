namespace Defyle.WebApi.Policy.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class GrantPolicyRequest
  {
    [Required]
    public string Subject { get; set; }
    
    public string Object { get; set; }

    [Required]
    public string Permission { get; set; }
  }
}