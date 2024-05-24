namespace Defyle.WebApi.Configuration
{
  using Core.Infrastructure.Poco;
  using Grpc.Core;
  using Microsoft.Extensions.DependencyInjection;
  using Proto;

  public static class GrpcExtensions
  {
    public static void AddGrpc(this IServiceCollection services, GrpcSettings coreServer)
    {
      Channel channel = new Channel(coreServer.Target, ChannelCredentials.Insecure);
      services.AddSingleton(channel);
      
      var workspaceProtoClient = new WorkspaceServiceProto.WorkspaceServiceProtoClient(channel);
      services.AddSingleton(workspaceProtoClient);

      var inodeProtoClient = new InodeServiceProto.InodeServiceProtoClient(channel);
      services.AddSingleton(inodeProtoClient);
      
      var userProtoClient = new UserServiceProto.UserServiceProtoClient(channel);
      services.AddSingleton(userProtoClient);
      
      var roleProtoClient = new RoleServiceProto.RoleServiceProtoClient(channel);
      services.AddSingleton(roleProtoClient);
      
      var fileProtoClient = new FileServiceProto.FileServiceProtoClient(channel);
      services.AddSingleton(fileProtoClient);
      
      var policyProtoClient = new PolicyServiceProto.PolicyServiceProtoClient(channel);
      services.AddSingleton(policyProtoClient);

      var modelProtoClient = new ModelServiceProto.ModelServiceProtoClient(channel);
      services.AddSingleton(modelProtoClient);
    }
  }
}