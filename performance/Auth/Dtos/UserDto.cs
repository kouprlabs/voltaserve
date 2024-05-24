namespace Defyle.WebApi.Auth.Dtos
{
  using System.Collections.Generic;

  public class UserDto
	{
    public string Id { get; set; }

    public string FullName { get; set; }

    public string Image { get; set; }

    public string Email { get; set; }

    public string Username { get; set; }

    public bool IsEmailConfirmed { get; set; }

    public bool IsSuperuser { get; set; }

    public bool IsSystem { get; set; }

    public bool IsLdap { get; set; }

    public IEnumerable<string> Permissions { get; set; }

    public long CreateTime { get; set; }
    
    public long UpdateTime { get; set; }
	}
}