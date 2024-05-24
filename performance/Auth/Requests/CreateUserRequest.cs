namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class CreateUserRequest
	{
    [Required]
    public string FullName { get; set; }

    public string Image { get; set; }
    
    [Required]
    public string Username { get; set; }
    
    [EmailAddress]
    [Required]
    public string Email { get; set; }

    [Required]
    public string Password { get; set; }

    public bool IsSuperuser { get; set; }

    public bool IsSystem { get; set; }
  }
}