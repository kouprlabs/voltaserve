namespace Defyle.WebApi.Inode.Controllers
{
  using System;
  using System.Collections.Generic;
  using System.IO;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Exceptions;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Infrastructure.Exceptions;
  using Core.Infrastructure.Poco;
  using Core.Inode.Models;
  using Core.Inode.Pocos;
  using Core.Inode.Services;
  using Core.Storage.Exceptions;
  using Core.Storage.Models;
  using Core.Storage.Services;
  using Core.Streaming.Services;
  using Core.Workspace.Models;
  using Core.Workspace.Services;
  using Dtos;
  using Filters;
  using Infrastructure.Dtos;
  using Infrastructure.Responses;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Http;
  using Microsoft.AspNetCore.Mvc;
  using Microsoft.Extensions.Logging;
  using Requests;
  using Responses;
  using Swashbuckle.AspNetCore.Annotations;
  using CopyInodesRequest = Requests.CopyInodesRequest;
  using CreateDirectoryRequest = Requests.CreateDirectoryRequest;
  using CreateFileRequest = Requests.CreateFileRequest;
  using DeleteInodesRequest = Requests.DeleteInodesRequest;
  using File = Core.Storage.Models.File;
  using InodeSearchRequest = Requests.InodeSearchRequest;
  using MoveInodesChildrenRequest = Requests.MoveInodesChildrenRequest;
  using MoveInodesRequest = Requests.MoveInodesRequest;
  using UpdateInodeRequest = Requests.UpdateInodeRequest;

  [Route("workspaces/{workspaceId}/inodes")]
	[Authorize]
	[PartitionIdCheck]
  [ApiExplorerSettings(GroupName = "Inodes")]
	public class InodesController : BaseInodeController
	{
    private readonly WorkspaceService _workspaceService;
		private readonly InodeService _service;
    private readonly UserService _userService;
    private readonly SftpService _sftpService;
    private readonly WebpStreamService _webpStreamService;
    private readonly PdfStreamService _pdfStreamService;
    private readonly OcrPdfStreamService _ocrPdfStreamService;
    private readonly OriginalStreamService _originalStreamService;
    private readonly FileService _fileService;
    private readonly IMapper _mapper;
    private readonly ILogger _logger;

    public InodesController(
      InodeService service,
      WorkspaceService workspaceService,
      UserService userService,
      SftpService sftpService,
      FileService fileService,
      WebpStreamService webpStreamService,
      PdfStreamService pdfStreamService,
      OcrPdfStreamService ocrPdfStreamService,
      OriginalStreamService originalStreamService,
      IMapper mapper,
      ILogger<InodesController> logger)
      : base(service, workspaceService, userService)
		{
      _workspaceService = workspaceService;
			_service = service;
      _userService = userService;
      _sftpService = sftpService;
      _webpStreamService = webpStreamService;
      _pdfStreamService = pdfStreamService;
      _ocrPdfStreamService = ocrPdfStreamService;
      _originalStreamService = originalStreamService;
      _fileService = fileService;
      _mapper = mapper;
      _logger = logger;
    }
    
    [HttpPost("createFile")]
    [SwaggerOperation("Create file", OperationId = "createFile")]
    [ProducesResponseType(typeof(InodeDto), 200)]
    public async Task<IActionResult> CreateFileAsync(IFormFile file, string workspaceId,
      [FromQuery] string parentId = "0",
      [FromQuery] bool indexContent = false, [FromQuery] bool plainTextError = false,
      [FromQuery] string password = null)
    {
      try
      {
        User user = await _userService.FindAsync(UserId);
        string effectiveParentId = await GetEffectiveNodeIdAsync(workspaceId, parentId);

        Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
        
        InodeFacet parent = await _service.FindOneAsync(effectiveParentId, user);

        var created = await _service.CreateFromFileAsync(workspace, parent, indexContent, password, file, user);

        InodeFacetDto dto = _mapper.Map<InodeFacetDto>(created);

        return Created(new Uri($"workspaces/{workspaceId}/inodes/getInformation/{dto.Id}", UriKind.Relative), dto);
      }
      catch (InternalServerErrorException e)
      {
        return new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 422};
      }
      catch (ResourceNotFoundException e)
      {
        if (plainTextError)
        {
          return NotFound(string.Join(' ', e.Errors.Select(error => error.Description)));
        }
        else
        {
          return NotFound(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))});
        }
      }
      catch (AuthorizationException e)
      {
        if (plainTextError)
        {
          return new ObjectResult(string.Join(' ', e.Errors.Select(error => error.Description)))
          {
            StatusCode = 403
          };
        }
        else
        {
          return new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 403};
        }
      }
      catch (StorageUsageExceededException e)
      {
        if (plainTextError)
        {
          return new ObjectResult(string.Join(' ', e.Errors.Select(error => error.Description)))
          {
            StatusCode = 403
          };
        }
        else
        {
          return new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 403};
        }
      }
    }
    
    [HttpPost("createDirectory")]
    [SwaggerOperation("Create directory", OperationId = "createDirectory")]
    [ProducesResponseType(typeof(InodeDto), 200)]
    public async Task<IActionResult> CreateDirectoryAsync(string workspaceId, [FromBody] CreateDirectoryRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }

      User user = await _userService.FindAsync(UserId);

      string effectiveParentId = await GetEffectiveNodeIdAsync(workspaceId, request.ParentId);
      request.ParentId = effectiveParentId;

      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      Inode created = await _service.CreateDirectoryAsync(workspace, request.ParentId, request.Name, user);
      InodeFacetDto dto = _mapper.Map<InodeFacetDto>(created);
      return Created(new Uri($"workspaces/{workspaceId}/inodes/getInformation/{dto.Id}", UriKind.Relative), dto);
    }

    [HttpPost("createEmptyFile")]
    [SwaggerOperation("Create empty file", OperationId = "createEmptyFile")]
    [ProducesResponseType(typeof(InodeDto), 200)]
    public async Task<IActionResult> CreateEmptyFileAsync(string workspaceId, [FromBody] CreateFileRequest request)
    {
      if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      User user = await _userService.FindAsync(UserId);

      string effectiveParentId = await GetEffectiveNodeIdAsync(workspaceId, request.ParentId);
      request.ParentId = effectiveParentId;
        
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      Inode created = await _service.CreateFileAsync(workspace, request.ParentId, request.Name , user);
      InodeFacetDto dto = _mapper.Map<InodeFacetDto>(created);
      
      return Created(new Uri($"workspaces/{workspaceId}/inodes/getInformation/{dto.Id}", UriKind.Relative), dto);
    }

    [HttpGet("{id}")]
    [SwaggerOperation("Get information", OperationId = "getInformation")]
		[ProducesResponseType(typeof(InodeDto), 200)]
		public async Task<IActionResult> GetInformationAsync(string workspaceId, string id)
		{
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Inode inode = await _service.FindOneAsync(effectiveNodeId, user);
      InodeFacetDto dto = _mapper.Map<InodeFacetDto>(inode);

      return Ok(dto);
		}

		[HttpPost("delete")]
    [SwaggerOperation("Delete", OperationId = "delete")]
		[ProducesResponseType(typeof(DeleteInodesResponse), 200)]
		public async Task<IActionResult> DeleteAsync(string workspaceId, [FromBody] DeleteInodesRequest request)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      List<string> effectiveIds = new List<string>();
      foreach (var id in request.Ids)
      {
        string effectiveId = await GetEffectiveNodeIdAsync(workspaceId, id); 
        effectiveIds.Add(effectiveId);
      }

      request.Ids = effectiveIds;

      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);

      var result = await _service.DeleteManyAsync(workspace, request.Ids, user);

      var response = _mapper.Map<DeleteInodesResponse>(result);
      if (response.Errors.Any())
      {
        return StatusCode(207, response);
      }
      else
      {
        return Ok(response);
      }
		}

		[HttpPost("{id}/move")]
    [SwaggerOperation("Move", OperationId = "move")]
		[ProducesResponseType(typeof(MoveInodesResponse), 200)]
		public async Task<IActionResult> MoveAsync(string workspaceId, string id,
      [FromBody] MoveInodesRequest request)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);

      var result = await _service.MoveManyAsync(workspace, inode, request.Ids, user);

      var response = _mapper.Map<MoveInodesResponse>(result);
      if (response.Errors.Any())
      {
        return StatusCode(207, response);
      }
      else
      {
        return Ok(response);
      }
		}

		[HttpPost("{id}/copy")]
    [SwaggerOperation("Copy", OperationId = "copy")]
		[ProducesResponseType(typeof(CopyInodesResponse), 200)]
		public async Task<IActionResult> CopyAsync(string workspaceId, string id,
      [FromBody] CopyInodesRequest request)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);

      CopyInodesResult result = await _service.CopyManyAsync(workspace, inode, request.Ids, user);

      var response = _mapper.Map<CopyInodesResponse>(result);
      if (response.Errors.Any())
      {
        return StatusCode(207, response);
      }
      else
      {
        return Ok(response);
      }
		}

		[HttpPatch("{id}")]
    [SwaggerOperation("Update", OperationId = "update")]
		[ProducesResponseType(typeof(InodeDto), 200)]
		public async Task<IActionResult> UpdateAsync(string workspaceId, string id,
			[FromBody] UpdateInodeRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);
      InodeFacet updated = await _service.UpdateNameAsync(workspace, inode, request.Name, user);

      InodeFacetDto dto = _mapper.Map<InodeFacetDto>(updated);

      return Ok(dto);
		}

		[HttpGet("{id}/children")]
    [SwaggerOperation("Get children", OperationId = "getChildren")]
		[ProducesResponseType(typeof(InodePagedResultDto), 200)]
		public async Task<IActionResult> GetChildrenAsync(string workspaceId, string id,
			[FromQuery] int page = 1, [FromQuery] int count = 50)
		{
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);
      InodePagedResult inodePagedResult = await _service.GetChildrenAsync(inode, page, count, user);

      InodePagedResultDto dto = _mapper.Map<InodePagedResultDto>(inodePagedResult);

      return Ok(dto);
		}

    [HttpPost("{id}/search")]
    [SwaggerOperation("Search", OperationId = "search")]
		[ProducesResponseType(typeof(InodePagedResultDto), 200)]
		public async Task<IActionResult> SearchAsync(string workspaceId, string id,
			[FromBody] InodeSearchRequest request, [FromQuery] int page = 0, [FromQuery] int count = 50)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodePagedResult pagedResult = await _service.SearchAsync(workspace,_mapper.Map<InodeSearchOptions>(request), page, count, user);

      InodeSearchResultResponse result = new InodeSearchResultResponse();
      result.Result = _mapper.Map<InodePagedResultDto>(pagedResult);
      result.Request = request;

      return Ok(result);
		}

    [HttpGet("{id}/children/getCount")]
    [SwaggerOperation("Get children count", OperationId = "getChildrenCount")]
		[ProducesResponseType(typeof(long), 200)]
		public async Task<IActionResult> GetChildrenCountAsync(string workspaceId, string id)
		{
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);

      long count = await _service.CountChildrenAsync(inode, user);

      return Ok(count);
		}

		[HttpPost("{id}/children/move")]
    [SwaggerOperation("Move children", OperationId = "moveChildren")]
		[ProducesResponseType(typeof(MoveInodesChildrenResponse), 200)]
		public async Task<IActionResult> MoveChildrenAsync(string workspaceId, string id,
      [FromBody] MoveInodesChildrenRequest request)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);

      var result = await _service.MoveManyChildrenAsync(workspace, inode, request.Ids, user);

      var response = _mapper.Map<MoveInodesChildrenResponse>(result);
      if (response.Errors.Any())
      {
        return StatusCode(207, response);
      }
      else
      {
        return Ok(response);
      }
		}

    [HttpPost("{id}/sftpImport")]
    [SwaggerOperation("SFTP import", OperationId = "sftpImport")]
		public async Task<IActionResult> SftpImportAsync(string workspaceId, string id,
			[FromBody] SftpImportRequest request, [FromQuery] string password = null)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
      
			User user = await _userService.FindAsync(UserId);
			Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
			Inode inode = await _service.FindOneAsync(effectiveNodeId, user);

			if (request.Port == 0)
			{
				request.Port = 22;
			}

			if (string.IsNullOrWhiteSpace(request.Directory))
			{
				request.Directory = "/";
			}

			// Ensure it's not a whitespace
			if (string.IsNullOrWhiteSpace(request.Username))
			{
				request.Username = null;
			}

			// Ensure it's not a whitespace
			if (string.IsNullOrWhiteSpace(request.Password))
			{
				request.Password = null;
			}

			await _sftpService.ImportAsync(user, workspace, inode, _mapper.Map<SftpImportOptions>(request), password);

			return Ok();
		}

		[HttpPost("{id}/upload")]
    [SwaggerOperation("Upload", OperationId = "upload")]
		[ProducesResponseType(typeof(void), 200)]
		public async Task<IActionResult> UploadAsync(string workspaceId, string id, IFormFile file,
			[FromQuery] string password = null, [FromQuery] bool plainTextError = false,
      [FromQuery] bool indexContent = false)
		{
			try
			{
				User user = await _userService.FindAsync(UserId);
				Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
				InodeFacet inode = await _service.FindOneAsync(id, user);

				await _service.UpdateFileAsync(workspace, inode, file.FileName, file.Length, indexContent, file.OpenReadStream(), password, user);

				return Ok();
			}
			catch (InternalServerErrorException e)
      {
        return plainTextError ? NotFound(string.Join(' ', e.Errors.Select(error => error.Description))) :
          new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 422};
      }
			catch (ResourceNotFoundException e)
      {
        return plainTextError ? NotFound(string.Join(' ', e.Errors.Select(error => error.Description))) :
          NotFound(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))});
      }
			catch (AuthorizationException e)
      {
        return plainTextError ? NotFound(string.Join(' ', e.Errors.Select(error => error.Description))) :
          new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 403};
      }
		}

		[HttpGet("{id}/downloadFile")]
    [SwaggerOperation("Download file", OperationId = "downloadFile")]
		[ProducesResponseType(typeof(FileResult), 200)]
		public async Task<IActionResult> DownloadFileAsync(string workspaceId, string id,
			[FromQuery] string password = null)
		{
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
        
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        throw new ResourceNotFoundException().WithError(Error.PreviewNotFoundError);
      }
      
      Stream stream = await _originalStreamService.GetStreamAsync(workspace, file, password);
      return File(stream, file.GetMime(), inode.Name);
		}
    
    [HttpGet("{id}/downloadPreview")]
    [SwaggerOperation("Download preview", OperationId = "downloadPreview")]
    [ProducesResponseType(typeof(FileResult), 200)]
    public async Task<IActionResult> DownloadPreviewAsync(string workspaceId, string id, [FromQuery] string password)
    {
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(id, user);

      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        throw new ResourceNotFoundException().WithError(Error.PreviewNotFoundError);
      }

      List<IStreamService> streamServices = new List<IStreamService> {_ocrPdfStreamService};

      string fileType = await _fileService.GetFileTypeAsync(file.Mime);
      string fileCategory = await _fileService.GetFileCategoryAsync(fileType);
      switch (fileCategory)
      {
        case "image":
          streamServices.Add(_webpStreamService);
          break;
        case "document":
          streamServices.Add(_pdfStreamService);
          break;
      }
      
      streamServices.Add(_originalStreamService);

      List<string> failedServices = new List<string>();
      List<string> failedLocalPaths = new List<string>();
      List<string> failedS3Keys = new List<string>();
      List<string> exceptions = new List<string>();
      foreach (IStreamService service in streamServices)
      {
        try
        {
          Stream stream = await service.GetStreamAsync(workspace, file, password);
          return File(stream, file.GetMime(), inode.Name);
        }
        catch (Exception e)
        {
          failedServices.Add(service.GetType().ToString());
          failedLocalPaths.Add(await service.GetLocalPathAsync(workspace, file));
          failedS3Keys.Add(await service.GetS3KeyAsync(workspace, file));
          exceptions.Add(e.Message);
        }
      }
      
      _logger.LogCritical( $"Failed to acquire stream for file {file.Id}, " +
                           $"category {fileCategory}, type {fileType}, mime {file.GetMime()}, " +
                           $"stream services [{string.Join(",", failedServices.ToArray())}], " +
                           $"local paths [{string.Join(",", failedLocalPaths.ToArray())}], " +
                           $"S3 keys [{string.Join(",", failedS3Keys.ToArray())}], " +
                           $"exceptions [{string.Join(",", exceptions.ToArray())}].");
      return NotFound();
    }
    
    [HttpGet("{id}/getStorageUsage")]
    [SwaggerOperation("Get storage usage", OperationId = "getStorageUsage")]
    [PartitionIdCheck]
    [ProducesResponseType(typeof(StorageUsage), 200)]
    public async Task<IActionResult> GetStorageUsageAsync(string workspaceId, string id)
    {
      string effectiveNodeId = await GetEffectiveNodeIdAsync(workspaceId, id);
      
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet inode = await _service.FindOneAsync(effectiveNodeId, user);

      StorageUsage storageUsage = await _service.GetStorageUsageAsync(workspace, inode, user);

      return Ok(storageUsage);
    }
  }
}