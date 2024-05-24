namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class PasswordSettings
	{
    [JsonProperty("digits")]
		public bool Digits { get; set; }

    [JsonProperty("requiredLength")]
		public int RequiredLength { get; set; }

    [JsonProperty("unique")]
		public bool Unique { get; set; }

    [JsonProperty("lowercase")]
		public bool Lowercase { get; set; }

    [JsonProperty("nonAlphanumeric")]
		public bool NonAlphanumeric { get; set; }

    [JsonProperty("uppercase")]
		public bool Uppercase { get; set; }
	}
}