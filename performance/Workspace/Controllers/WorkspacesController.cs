namespace Defyle.WebApi.Workspace.Controllers
{
  using System.Collections.Generic;
  using System.IO;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Workspace.Pocos;
  using Core.Workspace.Services;
  using Dtos;
  using Filters;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Http;
  using Microsoft.AspNetCore.Mvc;
  using Requests;
  using Responses;
  using Services;
  using Swashbuckle.AspNetCore.Annotations;
  using Workspace = Core.Workspace.Models.Workspace;

  [Route("workspaces")]
	[Authorize]
	[ProducesResponseType(typeof(WorkspaceDto), 200)]
  [ApiExplorerSettings(GroupName = "Workspaces")]
	public class WorkspacesController : BaseController
	{
    private readonly WorkspaceService _workspaceService;
    private readonly WorkspaceDtoService _workspaceDtoService;
		private readonly UserService _userService;
    private readonly IMapper _mapper;

    public WorkspacesController(
      WorkspaceService workspaceService,
			WorkspaceDtoService workspaceDtoService,
			UserService userService,
      IMapper mapper)
		{
      _workspaceService = workspaceService;
      _workspaceDtoService = workspaceDtoService;
			_userService = userService;
      _mapper = mapper;
    }

    [HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
		[ProducesResponseType(typeof(WorkspaceDto), 200)]
		public async Task<IActionResult> CreateAsync([FromBody] CreateWorkspaceRequest request)
		{
			if (!ModelState.IsValid)
			{
				return BadRequest(ModelState);
			}

      User user = await _userService.FindAsync(UserId);
      Workspace created = await _workspaceService.InsertAsync(_mapper.Map<CreateWorkspaceOptions>(request), user);
      WorkspaceDto dto = await _workspaceDtoService.CreateAsync(created, user);

      return Created("getWorkspaceInformation", dto);
		}
    
    [HttpGet]
    [SwaggerOperation("Find all", OperationId = "findAll")]
    public async Task<IEnumerable<WorkspaceDto>> FindAllAsync([FromQuery] bool images = false)
    {
      User user = await _userService.FindAsync(UserId);
      IEnumerable<Workspace> workspaces = await _workspaceService.FindAllAsync(user);
      var dtos = new List<WorkspaceDto>();

      foreach (Workspace workspace in workspaces)
      {
        WorkspaceDto dto = await _workspaceDtoService.CreateAsync(workspace, user);

        if (!images)
        {
          dto.Image = null;
        }

        dtos.Add(dto);
      }

      return dtos;
    }

		[HttpGet("{id}", Name = "getWorkspaceInformation")]
    [SwaggerOperation("Get information", OperationId = "getInformation")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(WorkspaceDto), 200)]
		public async Task<IActionResult> GetInformationAsync(string id)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);

      WorkspaceDto dto = await _workspaceDtoService.CreateAsync(workspace, user);

      return Ok(dto);
    }

		[HttpPatch("{id}")]
    [SwaggerOperation("Update", OperationId = "update")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(WorkspaceDto), 200)]
		public async Task<IActionResult> UpdateAsync(string id, [FromBody] UpdateWorkspaceRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);
      Workspace updated = await _workspaceService.UpdateNameAsync(workspace, request.Name, user);

      WorkspaceDto dto = await _workspaceDtoService.CreateAsync(updated, user);

      return Ok(dto);
		}

		[HttpPost("{id}/updateImage")]
    [SwaggerOperation("Update image", OperationId = "updateImage")]
		[ProducesResponseType(typeof(WorkspaceDto), 200)]
		public async Task<IActionResult> UpdateImageAsync(string id, IFormFile file)
		{
			string path = Path.GetTempFileName();

			using (var stream = new FileStream(path, FileMode.Create))
			{
				file.CopyTo(stream);
			}

      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);
      Workspace updated = await _workspaceService.UpdateImageAsync(workspace, path, user);
      WorkspaceDto dto = await _workspaceDtoService.CreateAsync(updated, user);

      return Ok(dto);
		}

		[HttpDelete("{id}")]
    [SwaggerOperation("Delete", OperationId = "delete")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(void), 200)]
		public async Task<IActionResult> DeleteAsync(string id)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);
      await _workspaceService.DeleteAsync(workspace, user);

      return Ok();
		}

    [HttpPost("{id}/verifyPassword")]
    [SwaggerOperation("Verify password", OperationId = "verifyPassword")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(void), 200)]
		public async Task<IActionResult> VerifyPasswordAsync(
			string id,
			[FromQuery] string password)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);
      _workspaceService.VerifyPassword(workspace, password);

      return Ok();
		}

		[HttpGet("{id}/getTransitKey")]
    [SwaggerOperation("Get transit key", OperationId = "getTransitKey")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(void), 200)]
		public async Task<IActionResult> GetTransitKeyAsync(string id)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);

      return Ok(new TransitKeyResponse
      {
        Key = workspace.TransitKey,
        Iv = workspace.TransitIv
      });
		}

		[HttpPost("{id}/updatePassword")]
    [SwaggerOperation("Update password", OperationId = "updatePassword")]
		[PartitionIdCheck]
		[ProducesResponseType(typeof(void), 200)]
		public async Task<IActionResult> UpdatePasswordAsync(string id, [FromBody] UpdateWorkspacePasswordRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(id, user);

      await _workspaceService.UpdatePasswordAsync(workspace, request.CurrentPassword, request.NewPassword, user);

      return Ok();
		}
	}
}