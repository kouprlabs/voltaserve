namespace Defyle.Core.Infrastructure.Configuration
{
  using Auth.Poco;
  using AutoMapper;
  using Inode.Models;
  using Inode.Pocos;
  using Policy.Pocos;
  using Proto;
  using Role.Models;
  using Role.Pocos;
  using Storage.Models;

  public class AutoMapperConfigurer
  {
    public void Configure(IMapperConfigurationExpression cfg)
    {
      cfg.CreateMap<string, string>().ConvertUsing(s => s ?? string.Empty);

      cfg.CreateMap<Core.Workspace.Models.Workspace, WorkspaceProto>();
      cfg.CreateMap<WorkspaceProto, Core.Workspace.Models.Workspace>()
        .ForMember(x => x.Tasks, opt => opt.Ignore());

      cfg.CreateMap<Core.Inode.Models.InodeFacet, InodeFacetProto>();
      cfg.CreateMap<InodeFacetProto, Core.Inode.Models.InodeFacet>();
      cfg.CreateMap<Core.Inode.Models.Inode, InodeProto>();
      cfg.CreateMap<InodeProto, Core.Inode.Models.Inode>();
      
      cfg.CreateMap<FileFacet, FileFacetProto>();
      cfg.CreateMap<FileFacetProto, FileFacet>();

      cfg.CreateMap<DirectoryFacet, DirectoryFacetProto>();
      cfg.CreateMap<DirectoryFacetProto, DirectoryFacet>();

      cfg.CreateMap<File, FileProto>();
      cfg.CreateMap<FileProto, File>();

      cfg.CreateMap<Core.Auth.Models.User, UserProto>();
      cfg.CreateMap<UserProto, Core.Auth.Models.User>();
      cfg.CreateMap<UserFindAllPagedResponseProto, UserPagedResult>();
      cfg.CreateMap<UserPagedResult, UserFindAllPagedResponseProto>();

      cfg.CreateMap<InodeGetChildrenPagedResponseProto, InodePagedResult>();
      cfg.CreateMap<InodePagedResult, InodeGetChildrenPagedResponseProto>();
        
      cfg.CreateMap<Core.Policy.Models.Policy, PolicyProto>();
      cfg.CreateMap<PolicyProto, Core.Policy.Models.Policy>();
      cfg.CreateMap<PolicyFindAllPagedResponseProto, PolicyPagedResult>();
      cfg.CreateMap<PolicyPagedResult, PolicyFindAllPagedResponseProto>();

      cfg.CreateMap<FileProperty, FilePropertyProto>();
      cfg.CreateMap<FilePropertyProto, FileProperty>();
      
      cfg.CreateMap<Role, RoleProto>();
      cfg.CreateMap<RoleProto, Role>();
      cfg.CreateMap<RoleFindAllPagedResponseProto, RolePagedResult>();
      cfg.CreateMap<RolePagedResult, RoleFindAllPagedResponseProto>();
    }
  }
}