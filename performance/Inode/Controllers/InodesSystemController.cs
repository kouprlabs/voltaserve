namespace Defyle.WebApi.Inode.Controllers
{
  using System;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Infrastructure.Poco;
  using Core.Inode.Models;
  using Core.Inode.Services;
  using Core.Workspace.Models;
  using Core.Workspace.Services;
  using Dtos;
  using Filters;
  using Microsoft.AspNetCore.Http;
  using Microsoft.AspNetCore.Mvc;
  using Requests;

  [Route("workspaces/{workspaceId}/inodes/system")]
  [PartitionIdCheck]
  [ApiExplorerSettings(GroupName = "Inodes")]
  public class InodesSystemController : BaseInodeController
  {
    private readonly CoreSettings _coreSettings;
    private readonly InodeService _service;
    private readonly WorkspaceService _workspaceService;
    private readonly UserService _userService;
    private readonly IMapper _mapper;

    public InodesSystemController(
      CoreSettings coreSettings,
      InodeService service,
      WorkspaceService workspaceService,
      UserService userService,
      IMapper mapper)
      : base(service, workspaceService, userService)
    {
      _coreSettings = coreSettings;
      _service = service;
      _workspaceService = workspaceService;
      _userService = userService;
      _mapper = mapper;
    }
    
    [HttpPost("createDirectory")]
    [ApiExplorerSettings(IgnoreApi = true)]
    [ProducesResponseType(typeof(InodeDto), 200)]
    public async Task<IActionResult> SystemCreateDirectoryAsync(string workspaceId, [FromQuery] string userId,
      [FromQuery] string parentId, [FromQuery] string name)
    {
      if (!_coreSettings.AllowedHosts.Contains(Request.Host.Host))
      {
        return NotFound();
      }

      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }

      User user = await _userService.FindAsync(userId);

      string effectiveParentId = await GetEffectiveNodeIdAsync(workspaceId, parentId);

      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      var request = new CreateDirectoryRequest
      {
        Name = name,
        ParentId = effectiveParentId
      };
      Inode created = await _service.CreateDirectoryAsync(workspace, request.ParentId, request.Name, user);

      return Created(new Uri($"workspaces/{workspaceId}/inodes/getInformation/{created.Id}", UriKind.Relative), created);
    }
    
    [HttpPost("createFile")]
    [ApiExplorerSettings(IgnoreApi = true)]
    [ProducesResponseType(typeof(void), 200)]
    public async Task<IActionResult> SystemCreateFileAsync(IFormFile file, string workspaceId, [FromQuery] string userId,
      [FromQuery] string parentNodeId = "0",
      [FromQuery] bool indexContent = false,
      [FromQuery] string password = null)
    {
      if (!_coreSettings.AllowedHosts.Contains(Request.Host.Host))
      {
        return NotFound();
      }

      User user = await _userService.FindAsync(userId);
      
      string effectiveParentId = await GetEffectiveNodeIdAsync(workspaceId, parentNodeId);

      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      InodeFacet parent = await _service.FindOneAsync(effectiveParentId, user);

      var created =  await _service.CreateFromFileAsync(workspace, parent, indexContent, password, file, user);

      return Created(new Uri($"workspaces/{workspaceId}/inodes/getInformation/{created.Id}", UriKind.Relative), created); 
    }
  }
}