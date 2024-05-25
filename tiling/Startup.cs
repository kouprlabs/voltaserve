namespace Voltaserve.Tiling
{
    using System.IO;
    using Voltaserve.Tiling.Services;
    using Microsoft.AspNetCore.Builder;
    using Microsoft.AspNetCore.Http.Features;
    using Microsoft.Extensions.Configuration;
    using Microsoft.Extensions.DependencyInjection;

    public class Startup
    {
        private IConfiguration _configuration;

        public Startup()
        {
            var outputDirectory = "out";
            if (!Directory.Exists(outputDirectory))
            {
                Directory.CreateDirectory(outputDirectory);
            }
        }

        private IConfiguration Configuration
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

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers()
              .AddNewtonsoftJson();

            services.Configure<FormOptions>(options =>
            {
                options.MultipartBodyLengthLimit = Configuration.GetValue<int>("multipartBodyLengthLimit");
            });

            services.AddWebEncoders();

            services.AddTransient<TilesService>();
        }

        public void Configure(IApplicationBuilder app)
        {
            app.UseRouting();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapDefaultControllerRoute();
            });
        }
    }
}