namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class CreateAccountRequest
	{
		[EmailAddress]
		[Required]
		public string Email { get; set; }

		[Required]
    public string Password { get; set; }

		[Required]
		public string FullName { get; set; }

		public string Image { get; set; }
	}
}