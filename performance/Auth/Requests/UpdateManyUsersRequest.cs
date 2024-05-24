namespace Defyle.WebApi.Auth.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class UpdateManyUsersRequest
  {
    [Required]
    public IEnumerable<UpdateUserRequest> Objects { get; set; }
  }
}