namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class ConfirmEmailRequest
	{
		[Required]
		public string Token { get; set; }
	}
}