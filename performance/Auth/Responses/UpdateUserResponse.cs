namespace Defyle.WebApi.Auth.Responses
{
  using System.Collections.Generic;
  using Core.Auth.Models;
  using Core.Infrastructure.Poco;

  public class UpdateUserResponse
  {
    public List<Error> Errors { get; set; } = new List<Error>();
    public User Result { get; set; }
  }
}