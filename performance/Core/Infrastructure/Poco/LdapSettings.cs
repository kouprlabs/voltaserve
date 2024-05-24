namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class LdapSettings
	{
    [JsonProperty("host")]
		public string Host { get; set; }

    [JsonProperty("port")]
		public int Port { get; set; }

    [JsonProperty("bindDn")]
		public string BindDn { get; set; }

    [JsonProperty("bindPassword")]
		public string BindPassword { get; set; }

    [JsonProperty("searchBase")]
		public string SearchBase { get; set; }

    [JsonProperty("searchFilter")]
		public string SearchFilter { get; set; }

    [JsonProperty("usernameAttribute")]
		public string UsernameAttribute { get; set; }

    [JsonProperty("adminCn")]
		public string AdminCn { get; set; }
	}
}