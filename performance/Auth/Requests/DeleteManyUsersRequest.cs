namespace Defyle.WebApi.Auth.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class DeleteManyUsersRequest
  {
    [Required]
    public IEnumerable<string> Ids { get; set; }
  }
}