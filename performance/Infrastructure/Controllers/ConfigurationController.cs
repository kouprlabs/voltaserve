namespace Defyle.WebApi.Infrastructure.Controllers
{
  using System.IO;
  using AutoMapper;
  using Core.Auth.Extensions;
  using Core.Infrastructure.Poco;
  using Dtos;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Responses;
  using SemVer;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("configuration")]
	[AllowAnonymous]
  [ApiExplorerSettings(GroupName = "Configuration")]
	public class ConfigurationController : Controller
	{
		private readonly CoreSettings _coreSettings;
    private readonly IMapper _mapper;

    public ConfigurationController(
			CoreSettings coreSettings,
      IMapper mapper)
    {
      _coreSettings = coreSettings;
      _mapper = mapper;
    }

		[HttpGet]
    [SwaggerOperation("Get configuration", OperationId = "getConfiguration")]
		[ProducesResponseType(typeof(string), 200)]
		public IActionResult GetConfiguration()
		{
			ConfigurationResponse response = new ConfigurationResponse();
			response.AuthenticationType = _coreSettings.AuthenticationType.NormalizedAuthenticationTypeString();

			Version version = new Version(
				System.IO.File.ReadAllText(Path.Combine(Directory.GetCurrentDirectory(), "config", "version.txt")));

      response.Version = new VersionDto
      {
        Build = version.Build,
        Major = version.Major,
        Minor = version.Minor,
        Patch = version.Patch,
        PreRelease = version.PreRelease
      };
			response.GravatarIntegration = _coreSettings.GravatarIntegration;

			return Ok(response);
		}
	}
}