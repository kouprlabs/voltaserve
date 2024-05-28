namespace Voltaserve.Mosaic
{
    using Voltaserve.Mosaic.Services;
    using Microsoft.AspNetCore.Builder;
    using Microsoft.AspNetCore.Http.Features;
    using Microsoft.Extensions.DependencyInjection;
    using dotenv.net.Utilities;
    using Minio;

    public class Startup
    {
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers();

            services.Configure<FormOptions>(options =>
            {
                options.MultipartBodyLengthLimit = EnvReader.GetIntValue("MULTIPART_BODY_LENGTH_LIMIT");
            });

            services.AddWebEncoders();

            services.AddMinio(configureClient => configureClient
                .WithEndpoint(EnvReader.GetStringValue("S3_URL"))
                .WithCredentials(EnvReader.GetStringValue("S3_ACCESS_KEY"), EnvReader.GetStringValue("S3_SECRET_KEY"))
                .WithRegion(EnvReader.GetStringValue("S3_REGION"))
                .WithSSL(EnvReader.GetBooleanValue("S3_SECURE"))
              );

            services.AddTransient<MosaicService>();
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