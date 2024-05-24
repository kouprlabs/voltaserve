namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class TokenSettings
	{
    [JsonProperty("accessTokenLifetime")]
		public int AccessTokenLifetime { get; set; }

    [JsonProperty("refreshTokenLifetime")]
		public int RefreshTokenLifetime { get; set; }

    [JsonProperty("tokenAudience")]
		public string TokenAudience { get; set; }

    [JsonProperty("tokenIssuer")]
		public string TokenIssuer { get; set; }
	}
}