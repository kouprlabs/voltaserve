namespace Defyle.WebApi
{
	using System.IO;
  using Core.Infrastructure.Poco;
  using Microsoft.AspNetCore;
	using Microsoft.AspNetCore.Hosting;
	using Microsoft.Extensions.Configuration;
	using Microsoft.Extensions.Logging;

  public class Program
	{
		private static IConfiguration _configuration;

		private static IConfiguration Configuration
		{
			get
			{
				if (_configuration == null)
				{
					_configuration = new ConfigurationBuilder()
						.SetBasePath(Directory.GetCurrentDirectory())
						.AddJsonFile(Path.Combine("config","appsettings.json"))
						.Build();
				}

				return _configuration;
			}
		}

		public static void Main(string[] args)
    {
      BuildWebHost(args).Run();
		}

		public static IWebHost BuildWebHost(string[] args)
		{
			CoreSettings coreSettings = Configuration.Get<CoreSettings>();
      
			return WebHost.CreateDefaultBuilder(args)
				.UseKestrel(options => { options.Limits.MaxRequestBodySize = coreSettings.MultipartBodyLengthLimit; })
				.UseStartup<Startup>()
				.UseUrls(coreSettings.Url)
				.ConfigureLogging((webHostBuilderContext, loggingBuilder) =>
				{
					loggingBuilder.ClearProviders();
          loggingBuilder.AddConsole();
        })
        .Build();
		}
	}
}