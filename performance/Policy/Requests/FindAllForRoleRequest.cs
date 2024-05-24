namespace Defyle.WebApi.Policy.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class FindAllForRoleRequest
  {
    [Required]
    public string RoleId { get; set; }

    [Required]
    public string Object { get; set; }
  }
}