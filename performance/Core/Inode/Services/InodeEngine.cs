namespace Defyle.Core.Inode.Services
{
  using System;
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using Auth.Models;
  using AutoMapper;
  using Infrastructure.Services;
  using Models;
  using Nest;
  using Pocos;
  using Proto;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class InodeEngine
  {
    private readonly IElasticClient _elasticClient;
    private readonly InodeServiceProto.InodeServiceProtoClient _protoClient;
    private readonly FileServiceProto.FileServiceProtoClient _fileProtoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;

    public InodeEngine(
      IElasticClient elasticClient,
      InodeServiceProto.InodeServiceProtoClient protoClient,
      FileServiceProto.FileServiceProtoClient fileProtoClient,
      GrpcExceptionTranslator exceptionTranslator,
      IMapper mapper)
    {
      _elasticClient = elasticClient;
      _protoClient = protoClient;
      _fileProtoClient = fileProtoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
    }

    public async Task<InodeFacet> CreateFileAsync(Workspace workspace, InodeFacet parent, string name, User user)
    {
      InodeFacet inode = new InodeFacet();
      inode.WorkspaceId = workspace.Id;
      inode.ParentId = parent.Id;
      inode.Name = name;
      inode.Type = Inode.InodeTypeFile;
      inode.CreateTime = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();

      var proto = await _protoClient.InsertAsync(new InodeInsertRequestProto
      {
        UserId = user.Id,
        Inode = _mapper.Map<InodeProto>(inode)
      });
      InodeFacet created = _mapper.Map<InodeFacet>(proto);

      return created;
    }

    public async Task<InodeFacet> CreateDirectoryAsync(Workspace workspace, InodeFacet parent, string name, User user)
    {
      InodeFacet inode = new InodeFacet();
      inode.WorkspaceId = workspace.Id;
      inode.ParentId = parent.Id;
      inode.Name = name;
      inode.Type = Inode.InodeTypeDirectory;
      inode.CreateTime = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();

      var proto = await _protoClient.InsertAsync(new InodeInsertRequestProto
      {
        UserId = user.Id,
        Inode = _mapper.Map<InodeProto>(inode)
      });

      return _mapper.Map<InodeFacet>(proto);
    }

    public async Task<InodeFacet> FindByIdAsync(string id, User user)
    {
      try
      {
        var proto = await _protoClient.FindAsync(new InodeFindRequestProto
        {
          UserId = user.Id,
          Id = id
        });

        return _mapper.Map<InodeFacet>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<InodeFacet> FindRootAsync(Workspace workspace, User user)
    {
      try
      {
        var proto = await _protoClient.FindRootAsync(new InodeFindRootRequestProto
        {
          UserId = user.Id,
          WorkspaceId = workspace.Id
        });
      
        return _mapper.Map<InodeFacet>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<string> CopyAsync(InodeFacet source, InodeFacet destination, User user)
    {
      try
      {
        var proto = await _protoClient.CopyAsync(new InodeCopyRequestProto
        {
          UserId = user.Id,
          Id = destination.Id,
          SourceId = source.Id
        });

        InodeFacet copy = _mapper.Map<InodeFacet>(proto);

        return copy.Id;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task MoveAsync(InodeFacet source, InodeFacet destination, User user)
    {
      try
      {
        await _protoClient.MoveAsync(new InodeMoveRequestProto
        {
          UserId = user.Id,
          Id = destination.Id,
          SourceId = source.Id
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task MoveChildrenAsync(InodeFacet source, InodeFacet destination, User user)
    {
      try
      {
        await _protoClient.MoveChildrenAsync(new InodeMoveChildrenRequestProto
        {
          UserId = user.Id,
          Id = destination.Id,
          SourceId = source.Id
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<string>> DeleteAsync(InodeFacet inode, User user)
    {
      try
      {
        var responseProto = await _fileProtoClient.FindAllForInodeAsync(new FileFindAllForInodeRequestProto
        {
          UserId = user.Id,
          NodeId = inode.Id
        });

        await _protoClient.DeleteAsync(new InodeDeleteRequestProto
        {
          Id = inode.Id,
          UserId = user.Id
        });

        return responseProto.Files.Select(p => p.Id).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task SetNameAsync(InodeFacet inode, string name, User user)
    {
      try
      {
        inode.Name = name;

        await _protoClient.UpdateAsync(new InodeUpdateRequestProto
        {
          UserId = user.Id,
          Inode = _mapper.Map<InodeProto>(inode)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task SetFileAsync(InodeFacet inode, File file, User user)
    {
      try
      {
        await _fileProtoClient.InsertAsync(new FileInsertRequestProto
        {
          UserId = user.Id,
          NodeId = inode.Id,
          File = _mapper.Map<FileProto>(file)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task SetTextAsync(InodeFacet inode, string text, User user)
    {
      try
      {
        inode.Text = text;
      
        await _protoClient.UpdateAsync(new InodeUpdateRequestProto
        {
          UserId = user.Id,
          Inode = _mapper.Map<InodeProto>(inode)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<List<Inode>> GetPathAsync(InodeFacet inode, User user)
    {
      try
      {
        var proto = await _protoClient.GetPathAsync(new InodeGetPathRequestProto
        {
          UserId = user.Id,
          Id = inode.Id
        });
        return proto.Inodes.Select(e => _mapper.Map<Inode>(e)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<InodePagedResult> GetChildrenAsync(InodeFacet inode, int page, int size, User user)
    {
      try
      {
        var proto = await _protoClient.GetChildrenPagedAsync(new InodeGetChildrenPagedRequestProto
        {
          UserId = user.Id,
          Id = inode.Id,
          Page = page,
          Size = size
        });

        return _mapper.Map<InodePagedResult>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<InodePagedResult> SearchAsync(Workspace workspace, int page, int size, InodeSearchOptions request, User user)
    {
      var searchCriteria = request.SearchCriteria.Trim();

      var nameQuery = GetSearchQuery(searchCriteria, "name");
      var textQuery = GetSearchQuery(searchCriteria, "text");

      string query =
        $"({nameQuery} OR {textQuery}) AND " +
        $"workspaceId:\"{workspace.Id}\"";

      if (!string.IsNullOrWhiteSpace(request.ParentId))
      {
        query += $" AND parentId:\"{request.ParentId}\"";
      }

      var response = await _elasticClient.SearchAsync<Inode>(e => e
        .Query(q => q
          .QueryString(c => c
            .Fields(f => f
              .Field("name")
              .Field("text")
              .Field("workspaceId")
              .Field("parentId"))
            .Query(query)))
        .From(0)
        .Size(1000) // Apply a hard limit of 1000
      );

      List<Inode> items = response.Hits.Select(e =>
      {
        Inode item = e.Source;
        item.Id = e.Id;
        return item;
      }).ToList();

      bool Filter(Inode e) => e.WorkspaceId == workspace.Id &&
                              (string.IsNullOrWhiteSpace(request.ParentId) || e.ParentId == request.ParentId) &&
                              (request.IncludeFiles != false || e.Type != Inode.InodeTypeFile) &&
                              (request.IncludeDirectories != false || e.Type != Inode.InodeTypeDirectory) &&
                              (request.CreatedAtFrom == null || e.CreateTime >= request.CreatedAtFrom.Value) &&
                              (request.CreatedAtTo == null || e.CreateTime <= request.CreatedAtTo.Value) &&
                              (request.UpdatedAtFrom == null || e.UpdateTime >= request.UpdatedAtFrom.Value) &&
                              (request.UpdatedAtTo == null || e.UpdateTime <= request.UpdatedAtTo.Value);
      
      // Apply in-memory filter
      List<Inode> filtered = items.Where(Filter).ToList();

      var proto = await _protoClient.FindManyAsync(new InodeFindManyRequestProto
      {
        UserId = user.Id,
        Ids = {filtered.Select(e => e.Id)}
      });
      
      var facets = proto.Inodes.Select(e => _mapper.Map<InodeFacet>(e)).ToList();

      // Apply in-memory paging
      IEnumerable<InodeFacet> paged = facets
        .Skip(page * size)
        .Take(size)
        .ToList();

      long totalElements = filtered.Count;
      long totalPages = totalElements / size;

      return new InodePagedResult(paged, totalPages, totalElements, page, size);
    }

    private static string GetSearchQuery(string criteria, string field)
    {
      string[] keywords = criteria.Split(" ");
      string keywordsQuery = string.Empty;
      foreach (string keyword in keywords)
      {
        keywordsQuery += $"{field}:{keyword}~ OR {field}:*{keyword}* OR {field}:/.*{keyword}.*/ ";
      }

      keywordsQuery = keywordsQuery.Trim();

      string query = keywordsQuery;
      if (keywords.Length > 1)
      {
        query +=
          $" OR {field}:\"{criteria}\"~ OR {field}:\"{criteria}\"~5 OR {field}:\"*{criteria}*\" OR {field}:/.*{criteria}.*/";
      }

      return query;
    }

    public async Task<long> CountChildrenAsync(InodeFacet inode, User user)
    {
      try
      {
        var proto = await _protoClient.CountChildrenAsync(new InodeCountChildrenRequestProto
        {
          UserId = user.Id,
          Id = inode.Id
        });
        return proto.Count;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<long> GetStorageUsageAsync(InodeFacet inode, User user)
    {
      try
      {
        var proto = await _protoClient.GetStorageUsageAsync(new InodeGetStorageUsageRequestProto
        {
          UserId = user.Id,
          Id = inode.Id
        });
        return proto.StorageUsage;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}