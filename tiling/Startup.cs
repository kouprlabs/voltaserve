namespace Voltaserve.Tiling
{
    using System.IO;
    using Voltaserve.Tiling.Services;
    using Microsoft.AspNetCore.Builder;
    using Microsoft.AspNetCore.Http.Features;
    using Microsoft.Extensions.DependencyInjection;
    using dotenv.net.Utilities;

    public class Startup
    {
        public Startup()
        {
            var outputDirectory = "out";
            if (!Directory.Exists(outputDirectory))
            {
                Directory.CreateDirectory(outputDirectory);
            }
        }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers()
              .AddNewtonsoftJson();

            services.Configure<FormOptions>(options =>
            {
                options.MultipartBodyLengthLimit = EnvReader.GetIntValue("MULTIPART_BODY_LENGTH_LIMIT");
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