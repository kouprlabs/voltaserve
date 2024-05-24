namespace Defyle.WebApi.Role.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class DeleteManyRolesRequest
  {
    [Required]
    public IEnumerable<string> Ids { get; set; }
  }
}