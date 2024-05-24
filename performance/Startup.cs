namespace Defyle.WebApi
{
  using System.IO;
  using AutoMapper;
  using Configuration;
  using Core.Auth.Services;
  using Core.Infrastructure.Poco;
  using Core.Infrastructure.Services;
  using Microsoft.AspNetCore.Builder;
  using Microsoft.AspNetCore.Hosting;
  using Microsoft.AspNetCore.Http.Features;
  using Microsoft.Extensions.Configuration;
	using Microsoft.Extensions.DependencyInjection;
  using Microsoft.Extensions.Logging;
  using Microsoft.OpenApi.Models;
  using Newtonsoft.Json;

  public class Startup
  {
    private readonly ILogger<Startup> _logger;
		private IConfiguration _configuration;
		private readonly CoreSettings _coreSettings;

    public Startup()
		{
      _coreSettings = Configuration.Get<CoreSettings>();
      
      _logger = LoggerFactory.Create(builder =>
        {
          builder.AddConfiguration(Configuration.GetSection("Logging"));
          builder.AddConsole();
        })
        .CreateLogger<Startup>();
      
      _logger.LogInformation(JsonConvert.SerializeObject(_coreSettings, Formatting.Indented));

			if (!Directory.Exists(_coreSettings.DataDirectory))
			{
        Directory.CreateDirectory(_coreSettings.DataDirectory);
			}
		}
    
    private IConfiguration Configuration
    {
      get
      {
        if (_configuration == null)
        {
          _configuration = new ConfigurationBuilder()
            .SetBasePath(Directory.GetCurrentDirectory())
            .AddJsonFile(Path.Combine("config", "appsettings.json"))
            .Build();
        }

        return _configuration;
      }
    }

		public void ConfigureServices(IServiceCollection services)
		{
      services.AddSingleton(_coreSettings);
      
      services.AddControllers()
        .AddNewtonsoftJson();
      
			services.Configure<FormOptions>(options =>
			{
				options.MultipartBodyLengthLimit = _coreSettings.MultipartBodyLengthLimit;
			});

			services.AddWebEncoders();
      
			services.AddCors();

			services.AddSecurity();

			services.AddMvc()
        .AddJsonOptions(options =>
        {
          options.JsonSerializerOptions.IgnoreNullValues = true;
        });

			services.AddSwaggerGen(c => {
        c.SwaggerDoc("v1", new OpenApiInfo {Title = "Defyle REST API", Version = "v1"});
        c.EnableAnnotations();
        
        c.DocInclusionPredicate((_, api) => !string.IsNullOrWhiteSpace(api.GroupName));
#pragma warning disable 618
        c.TagActionsBy(api => api.GroupName);
#pragma warning restore 618
      });

      services.AddSignalR();

      services.AddElasticsearch(_coreSettings.Elasticsearch, _logger);

      services.AddGrpc(_coreSettings.CoreService);
      
      services.AddAutoMapper();

			services.AddServices();
    }

    public void Configure(IApplicationBuilder app, IWebHostEnvironment env, IMapper mapper, UserService userService)
		{
      userService.CreateSystemUserIfNotExists();
      userService.CreateAdminUserIfNotExists();

      app.UseExceptionHandler(env, mapper);
      
      app.UseHsts();
      
      _coreSettings.AllowedCorsOrigins.Add(_coreSettings.WebClientUrl);

      app.UseRouting();
      app.UseCors(builder =>
      {
        builder.WithOrigins(_coreSettings.AllowedCorsOrigins.ToArray())
          .AllowAnyMethod()
          .AllowAnyHeader()
          .AllowCredentials();
      });

      app.UseAuthentication();
      app.UseAuthorization();

      app.UseEndpoints(endpoints =>
      {
        endpoints.MapDefaultControllerRoute();
        endpoints.MapHub<NotificationHub>("/notifications");
      });

      app.UseSwagger(c =>
      {
        c.RouteTemplate = "api-docs/{documentName}/swagger.json";
      });

      app.UseReDoc(c =>
      {
        c.RoutePrefix = "api-docs";
        c.SpecUrl = "v1/swagger.json";
      });

			app.UseRabbitMqListeners();
    }
	}
}