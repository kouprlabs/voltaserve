namespace Defyle.Core.Model.Services
{
  using System;
  using System.Collections.Generic;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Services;
  using Proto;

  public class ModelService
  {
    private readonly ModelServiceProto.ModelServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;

    public ModelService(
      ModelServiceProto.ModelServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator)
    {
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
    }

    public async Task<IEnumerable<string>> GetPermissionsAsync(string modelName)
    {
      try
      {
        var proto = await _protoClient.GetPermissionsAsync(new ModelGetPermissionsRequestProto{Name = modelName});
        return proto.Permissions;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<string> GetPropertyAsync(string modelName, string propertyName, User user)
    {
      try
      {
        var proto = await _protoClient.GetPropertyAsync(new ModelGetPropertyRequestProto
        {
          UserId = user.Id,
          ModelName = modelName,
          PropertyName = propertyName
        });
      
        return proto.Value;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}