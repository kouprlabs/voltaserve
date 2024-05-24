namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class ResetPasswordRequest
	{
		[Required]
		public string Token { get; set; }

		[Required]
		public string NewPassword { get; set; }
	}
}