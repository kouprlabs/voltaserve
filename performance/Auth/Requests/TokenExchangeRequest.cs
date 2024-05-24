namespace Defyle.WebApi.Auth.Requests
{
  using Microsoft.AspNetCore.Mvc;

  public class TokenExchangeRequest
	{
		[FromForm(Name = "grant_type")]
		public string GrantType { get; set; }

		[FromForm(Name = "username")]
		public string Username { get; set; }

		[FromForm(Name = "password")]
		public string Password { get; set; }

		[FromForm(Name = "refresh_token")]
		public string RefreshToken { get; set; }
	}
}