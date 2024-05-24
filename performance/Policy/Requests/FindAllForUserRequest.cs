namespace Defyle.WebApi.Policy.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class FindAllForUserRequest
  {
    [Required]
    public string Subject { get; set; }

    [Required]
    public string Object { get; set; }
  }
}