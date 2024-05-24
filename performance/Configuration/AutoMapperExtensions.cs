namespace Defyle.WebApi.Configuration
{
  using Auth.Dtos;
  using Auth.Requests;
  using Auth.Responses;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Poco;
  using Core.Infrastructure.Configuration;
  using Core.Infrastructure.Poco;
  using Core.Inode.Models;
  using Core.Inode.Pocos;
  using Core.Policy.Pocos;
  using Core.Role.Models;
  using Core.Role.Pocos;
  using Core.Storage.Models;
  using Core.Workspace.Models;
  using Core.Workspace.Pocos;
  using Infrastructure.Dtos;
  using Inode.Dtos;
  using Inode.Requests;
  using Inode.Responses;
  using Microsoft.Extensions.DependencyInjection;
  using Policy.Dtos;
  using Policy.Requests;
  using Role.Dtos;
  using Role.Requests;
  using Workspace.Dtos;
  using Workspace.Requests;

  public static class AutoMapperExtensions
  {
    public static void AddAutoMapper(this IServiceCollection services)
    {
      var configuration = new MapperConfiguration(cfg =>
      {
        new AutoMapperConfigurer().Configure(cfg);

        cfg.CreateMap<FileProperty, FilePropertyDto>();
        cfg.CreateMap<FilePropertyDto, FileProperty>();

        cfg.CreateMap<Inode, InodeDto>();
        cfg.CreateMap<FileFacet, FileFacetDto>();
        cfg.CreateMap<DirectoryFacet, DirectoryFacetDto>();
        cfg.CreateMap<InodeFacet, InodeFacetDto>();
        cfg.CreateMap<InodePagedResult, InodePagedResultDto>();
        cfg.CreateMap<CopyInodesResult, CopyInodesResponse>();
        cfg.CreateMap<DeleteInodesResult, DeleteInodesResponse>();
        cfg.CreateMap<MoveInodesResult, MoveInodesResponse>();
        cfg.CreateMap<MoveInodesChildrenResult, MoveInodesChildrenResponse>();

        cfg.CreateMap<CreateWorkspaceRequest, CreateWorkspaceOptions>();
        cfg.CreateMap<WorkspaceTask, WorkspaceTaskDto>();
        
        cfg.CreateMap<TokenExchangeRequest, TokenExchangeOptions>();
        
        cfg.CreateMap<UpdateUserResult, UpdateUserResponse>();

        cfg.CreateMap<SftpImportRequest, SftpImportRequest>();
        
        cfg.CreateMap<Error, ErrorDto>();
        cfg.CreateMap<User, UserDto>()
          .ForMember(x => x.Permissions, opts => opts.Ignore());
        cfg.CreateMap<Token, TokenDto>();
        cfg.CreateMap<UserPagedResult, UserPagedResultDto>();
        cfg.CreateMap<CreateUserRequest, User>()
          .ForMember(x => x.Id, opts => opts.Ignore())
          .ForMember(x => x.PasswordHash, opts => opts.MapFrom(p => p.Password))
          .ForMember(x => x.RefreshTokenValue, opts => opts.Ignore())
          .ForMember(x => x.RefreshTokenValidTo, opts => opts.Ignore())
          .ForMember(x => x.ResetPasswordToken, opts => opts.Ignore())
          .ForMember(x => x.EmailConfirmationToken, opts => opts.Ignore())
          .ForMember(x => x.IsEmailConfirmed, opts => opts.Ignore())
          .ForMember(x => x.IsLdap, opts => opts.Ignore())
          .ForMember(x => x.RoleId, opts => opts.Ignore())
          .ForMember(x => x.CreateTime, opts => opts.Ignore())
          .ForMember(x => x.UpdateTime, opts => opts.Ignore());
        cfg.CreateMap<UpdateUserRequest, User>()
          .ForMember(x => x.RefreshTokenValue, opts => opts.Ignore())
          .ForMember(x => x.RefreshTokenValidTo, opts => opts.Ignore())
          .ForMember(x => x.ResetPasswordToken, opts => opts.Ignore())
          .ForMember(x => x.EmailConfirmationToken, opts => opts.Ignore())
          .ForMember(x => x.IsEmailConfirmed, opts => opts.Ignore())
          .ForMember(x => x.IsLdap, opts => opts.Ignore())
          .ForMember(x => x.RoleId, opts => opts.Ignore())
          .ForMember(x => x.CreateTime, opts => opts.Ignore())
          .ForMember(x => x.UpdateTime, opts => opts.Ignore());
        cfg.CreateMap<CreateAccountRequest, User>()
          .ForMember(x => x.Id, opts => opts.Ignore())
          .ForMember(x => x.Username, opts => opts.Ignore())
          .ForMember(x => x.PasswordHash, opts => opts.MapFrom(p => p.Password))
          .ForMember(x => x.RefreshTokenValue, opts => opts.Ignore())
          .ForMember(x => x.RefreshTokenValidTo, opts => opts.Ignore())
          .ForMember(x => x.ResetPasswordToken, opts => opts.Ignore())
          .ForMember(x => x.EmailConfirmationToken, opts => opts.Ignore())
          .ForMember(x => x.IsEmailConfirmed, opts => opts.Ignore())
          .ForMember(x => x.IsSuperuser, opts => opts.Ignore())
          .ForMember(x => x.IsSystem, opts => opts.Ignore())
          .ForMember(x => x.IsLdap, opts => opts.Ignore())
          .ForMember(x => x.RoleId, opts => opts.Ignore())
          .ForMember(x => x.CreateTime, opts => opts.Ignore())
          .ForMember(x => x.UpdateTime, opts => opts.Ignore());

        cfg.CreateMap<InodeSearchRequest, InodeSearchOptions>();

        cfg.CreateMap<Core.Policy.Models.Policy, PolicyDto>();
        cfg.CreateMap<PolicyPagedResult, PolicyPagedResultDto>();
        cfg.CreateMap<CreatePolicyRequest, Core.Policy.Models.Policy>()
          .ForMember(x => x.Id, opts => opts.Ignore())
          .ForMember(x => x.CreateTime, opts => opts.Ignore());

        cfg.CreateMap<Role, RoleDto>();
        cfg.CreateMap<RolePagedResult, RolePagedResultDto>();
        cfg.CreateMap<CreateRoleRequest, Role>()
          .ForMember(x => x.Id, opts => opts.Ignore())
          .ForMember(x => x.CreateTime, opts => opts.Ignore())
          .ForMember(x => x.UpdateTime, opts => opts.Ignore());
        cfg.CreateMap<UpdateRoleRequest, Role>()
          .ForMember(x => x.CreateTime, opts => opts.Ignore())
          .ForMember(x => x.UpdateTime, opts => opts.Ignore());
      });
      configuration.AssertConfigurationIsValid();
      var mapper = configuration.CreateMapper();
      services.AddSingleton(mapper);
    }
  }
}