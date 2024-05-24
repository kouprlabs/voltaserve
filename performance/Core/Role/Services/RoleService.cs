namespace Defyle.Core.Role.Services
{
  using System;
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using Auth.Models;
  using AutoMapper;
  using Infrastructure.Services;
  using Models;
  using Pocos;
  using Proto;

  public class RoleService
  {
    private readonly RoleServiceProto.RoleServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;

    public RoleService(
      RoleServiceProto.RoleServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator,
      IMapper mapper)
    {
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
    }

    public async Task<Role> InsertAsync(Role role, User user)
    {
      try
      {
        var proto = await _protoClient.InsertAsync(new RoleInsertRequestProto
        {
          UserId = user.Id,
          Role = _mapper.Map<RoleProto>(role)
        });
      
        return _mapper.Map<Role>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<Role> FindAsync(string id, User user)
    {
      try
      {
        var proto = await _protoClient.FindAsync(new RoleFindRequestProto
        {
          UserId = user.Id,
          Id = id
        });

        return _mapper.Map<Role>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task DeleteAsync(Role role, User user)
    {
      try
      {
        await _protoClient.DeleteAsync(new RoleDeleteRequestProto
        {
          UserId = user.Id,
          Role = _mapper.Map<RoleProto>(role)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<Role> UpdateAsync(Role role, User user)
    {
      try
      {
        var proto = await _protoClient.UpdateAsync(new RoleUpdateRequestProto
        {
          UserId = user.Id,
          Role = _mapper.Map<RoleProto>(role)
        });
      
        return _mapper.Map<Role>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<Role>> FindAllAsync(User user)
    {
      try
      {
        var proto = await _protoClient.FindAllAsync(new RoleFindAllRequestProto
        {
          UserId = user.Id
        });

        return proto.Roles.Select(p => _mapper.Map<Role>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<RolePagedResult> FindAllPagedAsync(User user, int page, int size)
    {
      try
      {
        var proto = await _protoClient.FindAllPagedAsync(new RoleFindAllPagedRequestProto
        {
          UserId = user.Id,
          Page = page,
          Size = size
        });

        return _mapper.Map<RolePagedResult>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<Role>> FindAllForObjectAsync(string obj, User user)
    {
      try
      {
        var proto = await _protoClient.FindAllForObjectAsync(new RoleFindAllForObjectRequestProto
        {
          UserId = user.Id,
          Object = obj
        });

        return proto.Roles.Select(p => _mapper.Map<Role>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<Role>> FindAllWithoutPermissionsForObjectAsync(string obj, User user)
    {
      try
      {
        var proto = await _protoClient.FindAllWithoutPermissionsForObjectAsync(new RoleFindAllWithoutPermissionsForObjectRequestProto
        {
          UserId = user.Id,
          Object = obj
        });

        return proto.Roles.Select(p => _mapper.Map<Role>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}