namespace Defyle.WebApi.Role.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class UpdateRoleRequest
  {
    [Required]
    public string Id { get; set; }
    
    [Required]
    public string Name { get; set; }
  }
}