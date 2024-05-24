namespace Defyle.Core.Auth.Models
{
  public class User
	{
    public string Id { get; set; }

		public string FullName { get; set; }

		public string Image { get; set; }

		public string Email { get; set; }

		public string Username { get; set; }

		public string PasswordHash { get; set; }

    public string RefreshTokenValue { get; set; }
    
    public long RefreshTokenValidTo { get; set; }

    public string ResetPasswordToken { get; set; }

		public string EmailConfirmationToken { get; set; }

		public bool IsEmailConfirmed { get; set; }

		public bool IsSuperuser { get; set; }

    public bool IsSystem { get; set; }

		public bool IsLdap { get; set; }

    public string RoleId { get; set; }

    public long CreateTime { get; set; }
    
    public long UpdateTime { get; set; }
	}
}