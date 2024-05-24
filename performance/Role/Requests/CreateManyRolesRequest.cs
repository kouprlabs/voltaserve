namespace Defyle.WebApi.Role.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class CreateManyRolesRequest
  {
    [Required]
    public IEnumerable<CreateRoleRequest> Objects { get; set; }
  }
}