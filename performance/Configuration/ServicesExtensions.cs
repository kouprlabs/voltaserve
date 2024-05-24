namespace Defyle.WebApi.Configuration
{
  using Auth.Services;
  using Core.Auth.Services;
  using Core.Infrastructure.Services;
  using Core.Inode.Services;
  using Core.Model.Services;
  using Core.Ocr.Queue;
  using Core.Ocr.Services;
  using Core.Policy.Services;
  using Core.Preview.Queue;
  using Core.Preview.Services;
  using Core.Role.Services;
  using Core.Storage.Services;
  using Core.Streaming.Services;
  using Core.Workspace.Services;
  using Microsoft.Extensions.DependencyInjection;
  using Workspace.Services;

  public static class ServicesExtensions
  {
    public static void AddServices(this IServiceCollection services)
    {
      services.AddTransient<GrpcExceptionTranslator>();
      services.AddTransient<EmailValidationService>();
      services.AddTransient<Base64ImageService>();
      services.AddTransient<UserService>();
      services.AddTransient<UserDtoService>();
      services.AddTransient<PathService>();
      services.AddTransient<TileMapService>();
      services.AddTransient<TextExtractionService>();
      services.AddTransient<WebpStreamService>();
      services.AddTransient<PdfStreamService>();
      services.AddTransient<OcrPdfStreamService>();
      services.AddTransient<OriginalStreamService>();
      services.AddTransient<StorageService>();
      services.AddTransient<EncryptionService>();
      services.AddTransient<WorkspaceDtoService>();
      services.AddTransient<WorkspaceService>();
      services.AddTransient<InodeService>();
      services.AddTransient<InodeEngine>();
      services.AddTransient<InodeNotificationService>();
      services.AddTransient<PasswordService>();
      services.AddTransient<TokenService>();
      services.AddTransient<LdapService>();
      services.AddTransient<AccountService>();
      services.AddTransient<WorkspaceNotificationService>();
      services.AddTransient<QueueService>();
      services.AddTransient<JsonService>();
      services.AddTransient<PreviewQueueService>();
      services.AddTransient<OcrQueueService>();
      services.AddTransient<EmailQueueService>();
      services.AddTransient<PolicyService>();
      services.AddTransient<FileService>();
      services.AddTransient<RoleService>();
      services.AddTransient<ModelService>();
      
      services.AddSingleton<Dispatcher>();
      services.AddSingleton<SftpService>();
    }
  }
}