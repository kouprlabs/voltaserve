namespace Defyle.Core.Auth.Poco
{
  public class TokenExchangeOptions
	{
		public string GrantType { get; set; }

		public string Username { get; set; }
    
		public string Password { get; set; }
    
		public string RefreshToken { get; set; }
	}
}