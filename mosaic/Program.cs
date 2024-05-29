namespace Voltaserve.Mosaic
{
    using Microsoft.AspNetCore;
    using Microsoft.AspNetCore.Hosting;
    using Microsoft.Extensions.Logging;
    using dotenv.net.Utilities;
    using dotenv.net;

    public class Program
    {
        public static void Main(string[] args)
        {
            DotEnv.Load(new(overwriteExistingVars: false));
            BuildWebHost(args).Run();
        }

        public static IWebHost BuildWebHost(string[] args)
        {
            var url = EnvReader.GetStringValue("URL");
            var multipartBodyLengthLimit = EnvReader.GetIntValue("MULTIPART_BODY_LENGTH_LIMIT");

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