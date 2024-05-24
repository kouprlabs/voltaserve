namespace Defyle.WebApi.Role.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class UpdateManyRolesRequest
  {
    [Required]
    public IEnumerable<UpdateRoleRequest> Objects { get; set; }
  }
}