namespace Voltaserve.Tiling
{
    using System.IO;
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
                _configuration ??= new ConfigurationBuilder()
                    .SetBasePath(Directory.GetCurrentDirectory())
                    .AddJsonFile("appsettings.json")
                    .Build();
                return _configuration;
            }
        }

        public static void Main(string[] args)
        {
            BuildWebHost(args).Run();
        }

        public static IWebHost BuildWebHost(string[] args)
        {
            var url = Configuration.GetValue<string>("url");
            var multipartBodyLengthLimit = _configuration.GetValue<int>("multipartBodyLengthLimit");

            return WebHost.CreateDefaultBuilder(args)
              .UseKestrel(options => { options.Limits.MaxRequestBodySize = multipartBodyLengthLimit; })
              .UseStartup<Startup>()
              .UseUrls(url)
              .ConfigureLogging((webHostBuilderContext, loggingBuilder) =>
              {
                  loggingBuilder.ClearProviders();
                  loggingBuilder.AddConsole();
              })
              .Build();
        }
    }
}