namespace Defyle.WebApi.Auth.Requests
{
	public class UpdateUserRequest
	{
    public string Id { get; set; }
    
    public string FullName { get; set; }

    public string Image { get; set; }

    public string Email { get; set; }

    public string Username { get; set; }

    public string PasswordHash { get; set; }

    public bool IsEmailConfirmed { get; set; }

    public bool IsSuperuser { get; set; }

    public bool IsSystem { get; set; }

    public bool IsLdap { get; set; }
	}
}