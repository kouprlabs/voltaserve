namespace Defyle.Core.Storage.Services
{
  using System;
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using Auth.Models;
  using AutoMapper;
  using Infrastructure.Exceptions;
  using Infrastructure.Services;
  using Inode.Models;
  using Models;
  using Proto;

  public class FileService
  {
    private readonly FileServiceProto.FileServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;

    public FileService(
      FileServiceProto.FileServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator,
      IMapper mapper)
    {
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
    }

    public async Task<IEnumerable<File>> FindAllForInodeAsync(Inode inode, User user)
    {
      try
      {
        var proto = await _protoClient.FindAllForInodeAsync(new FileFindAllForInodeRequestProto
        {
          UserId = user.Id,
          NodeId = inode.Id
        });
        return proto.Files.Select(e => _mapper.Map<File>(proto));
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<File> FindLatestForInodeOrNullAsync(Inode inode, User user)
    {
      try
      {
        return await FindLatestForInodeAsync(inode, user);
      }
      catch (ResourceNotFoundException)
      {
        return null;
      }
    }
    
    public async Task<File> FindLatestForInodeAsync(Inode inode, User user)
    {
      try
      {
        var proto = await _protoClient.FindLatestForInodeAsync(new FileFindLatestForInodeRequestProto
        {
          UserId = user.Id,
          NodeId = inode.Id
        });
        return _mapper.Map<File>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task SetPropertyAsync(File file, string statusName, string value, User user)
    {
      try
      {
        await _protoClient.SetPropertyAsync(new FileSetPropertyRequestProto
        {
          UserId = user.Id,
          File = _mapper.Map<FileProto>(file),
          PropertyName = statusName,
          Value = value
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<string> GetPropertyAsync(File file, string statusName, User user)
    {
      try
      {
        var proto = await _protoClient.GetPropertyAsync(new FileGetPropertyRequestProto
        {
          UserId = user.Id,
          File = _mapper.Map<FileProto>(file),
          PropertyName = statusName
        });
        return proto.Value;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<FileProperty>> GetPropertiesAsync(File file, User user)
    {
      try
      {
        var proto = await _protoClient.GetPropertiesAsync(new FileGetPropertiesRequestProto
        {
          UserId = user.Id,
          File = _mapper.Map<FileProto>(file)
        });
        return proto.Properties.Select(e => _mapper.Map<FileProperty>(e));
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<string> GetFileTypeAsync(string mime)
    {
      try
      {
        var proto = await _protoClient.GetFileTypeAsync(new FileGetFileTypeRequestProto
        {
          Mime = mime
        });
        return proto.FileType;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<string> GetFileCategoryAsync(string fileType)
    {
      try
      {
        var proto = await _protoClient.GetFileCategoryAsync(new FileGetFileCategoryRequestProto
        {
          FileType = fileType
        });
        return proto.FileCategory;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}