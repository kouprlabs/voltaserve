namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class SendResetPasswordEmailRequest
	{
		[Required]
		public string Email { get; set; }
	}
}