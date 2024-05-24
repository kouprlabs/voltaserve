namespace Defyle.Core.Auth.Poco
{
  using System.Collections.Generic;
  using Infrastructure.Poco;
  using Models;

  public class UpdateUserResult
  {
    public List<Error> Errors { get; set; } = new List<Error>();
    
    public User Result { get; set; }
  }
}