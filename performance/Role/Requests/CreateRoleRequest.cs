namespace Defyle.WebApi.Role.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class CreateRoleRequest
  {
    [Required]
    public string Name { get; set; }
  }
}