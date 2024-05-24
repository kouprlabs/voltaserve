namespace Defyle.WebApi.Infrastructure.Responses
{
  using Dtos;

  public class ConfigurationResponse
	{
		public string AuthenticationType { get; set; }

		public VersionDto Version { get; set; }

		public bool GravatarIntegration { get; set; }
	}
}