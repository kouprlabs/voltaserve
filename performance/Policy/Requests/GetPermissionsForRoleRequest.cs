namespace Defyle.WebApi.Policy.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class GetPermissionsForRoleRequest
  {
    [Required]
    public string RoleId { get; set; }

    [Required]
    public string Object { get; set; }
  }
}