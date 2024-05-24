namespace Defyle.WebApi.Policy.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class RevokePolicyRequest
  {
    [Required]
    public string Subject { get; set; }

    [Required]
    public string Object { get; set; }

    [Required]
    public string Permission { get; set; }
  }
}