namespace Defyle.WebApi.Auth.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class CreateManyUsersRequest
  {
    [Required]
    public IEnumerable<CreateUserRequest> Objects { get; set; }
  }
}