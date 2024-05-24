namespace Defyle.Core.Policy.Services
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

  public class PolicyService
  {
    private readonly PolicyServiceProto.PolicyServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;

    public PolicyService(
      PolicyServiceProto.PolicyServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator,
      IMapper mapper)
    {
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
    }
    
    public async Task<Policy> InsertAsync(Policy role, User user)
    {
      try
      {
        var proto = await _protoClient.InsertAsync(new PolicyInsertRequestProto
        {
          UserId = user.Id,
          Policy = _mapper.Map<PolicyProto>(role)
        });
      
        return _mapper.Map<Policy>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task DeleteAsync(Policy role, User user)
    {
      try
      {
        await _protoClient.DeleteAsync(new PolicyDeleteRequestProto
        {
          UserId = user.Id,
          Policy = _mapper.Map<PolicyProto>(role)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<Policy> FindAsync(string id, User user)
    {
      try
      {
        var proto = await _protoClient.FindAsync(new PolicyFindRequestProto
        {
          UserId = user.Id,
          Id = id
        });

        return _mapper.Map<Policy>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<Policy>> FindAllAsync(User user)
    {
      try
      {
        var proto = await _protoClient.FindAllAsync(new PolicyFindAllRequestProto
        {
          UserId = user.Id
        });

        return proto.Policies.Select(p => _mapper.Map<Policy>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<PolicyPagedResult> FindAllPagedAsync(User user, int page, int size)
    {
      try
      {
        var proto = await _protoClient.FindAllPagedAsync(new PolicyFindAllPagedRequestProto
        {
          UserId = user.Id,
          Page = page,
          Size = size
        });

        return _mapper.Map<PolicyPagedResult>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task GrantAsync(string sub, string obj, string permission, User user)
    {
      try
      {
        await _protoClient.GrantAsync(new PolicyGrantRequestProto
        {
          UserId = user.Id,
          Object = obj ?? "",
          Subject = sub,
          Permission = permission
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task RevokeAsync(string sub, string obj, string permission, User user)
    {
      try
      {
        await _protoClient.RevokeAsync(new PolicyRevokeRequestProto
        {
          UserId = user.Id,
          Object = obj ?? "",
          Subject = sub,
          Permission = permission
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<Policy>> FindAllForUserAsync(string sub, string obj, User user)
    {
      try
      {
        var proto = await _protoClient.FindAllForUserAsync(new PolicyFindAllForUserRequestProto
        {
          UserId = user.Id,
          SubjectId = sub,
          Object = obj
        });
      
        return proto.Policies.Select(p => _mapper.Map<Policy>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<Policy>> FindAllForRoleAsync(string roleId, string obj, User user)
    {
      try
      {
        var proto = await _protoClient.FindAllForRoleAsync(new PolicyFindAllForRoleRequestProto
        {
          UserId = user.Id,
          RoleId = roleId,
          Object = obj
        });
      
        return proto.Policies.Select(p => _mapper.Map<Policy>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<string>> GetModelPermissionsForUserAsync(string subjectId, string model, User user)
    {
      try
      {
        var proto = await _protoClient.GetModelPermissionsForUserAsync(new PolicyGetModelPermissionsForUserRequestProto
        {
          UserId = user.Id,
          SubjectId = subjectId,
          Model = model
        });
      
        return proto.Permissions;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public IEnumerable<string> GetModelPermissionsForUser(string subjectId, string model, User user)
    {
      try
      {
        var proto = _protoClient.GetModelPermissionsForUser(new PolicyGetModelPermissionsForUserRequestProto
        {
          UserId = user.Id,
          SubjectId = subjectId,
          Model = model
        });
      
        return proto.Permissions;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<string>> GetObjectPermissionsForUserAsync(string subjectId, string obj, User user)
    {
      try
      {
        var proto = await _protoClient.GetObjectPermissionsForUserAsync(new PolicyGetObjectPermissionsForUserRequestProto
        {
          UserId = user.Id,
          SubjectId = subjectId,
          Object = obj
        });
      
        return proto.Permissions;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<string>> GetObjectPermissionsForRoleAsync(string roleId, string obj, User user)
    {
      try
      {
        var proto = await _protoClient.GetObjectPermissionsForRoleAsync(new PolicyGetObjectPermissionsForRoleRequestProto
        {
          UserId = user.Id,
          RoleId = roleId,
          Object = obj
        });
      
        return proto.Permissions;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}