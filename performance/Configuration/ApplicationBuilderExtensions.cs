namespace Defyle.WebApi.Configuration
{
  using Core.Inode.Services;
  using Core.Preview.Services;
  using Grpc.Core;
  using Microsoft.AspNetCore.Builder;
  using Microsoft.Extensions.DependencyInjection;
  using Microsoft.Extensions.Hosting;

  public static class ApplicationBuilderExtensions
	{
		private static Dispatcher Dispatcher { get; set; }
		private static SftpService SftpService { get; set; }

    private static Channel GrpcChannel { get; set; }

		public static IApplicationBuilder UseRabbitMqListeners(this IApplicationBuilder app)
		{
			Dispatcher = app.ApplicationServices.GetService<Dispatcher>();
			SftpService = app.ApplicationServices.GetService<SftpService>();
      GrpcChannel = app.ApplicationServices.GetService<Channel>();

			var life = app.ApplicationServices.GetService<IHostApplicationLifetime>();

			life.ApplicationStarted.Register(OnStarted);
			life.ApplicationStopping.Register(OnStopping);

			return app;
		}

		private static void OnStarted()
		{
			Dispatcher.Register();
			SftpService.Register();
		}

		private static void OnStopping()
		{
			Dispatcher.Deregister();
			SftpService.Deregister();
      GrpcChannel.ShutdownAsync().Wait();
		}
	}
}