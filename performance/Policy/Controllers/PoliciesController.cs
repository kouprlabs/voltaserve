namespace Defyle.WebApi.Policy.Controllers
{
  using System;
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Inode.Models;
  using Core.Inode.Services;
  using Core.Policy.Models;
  using Core.Policy.Pocos;
  using Core.Policy.Services;
  using Core.Workspace.Services;
  using Dtos;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Requests;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("policies")]
  [Authorize]
  [ApiExplorerSettings(GroupName = "Policies")]
  public class PoliciesController : BaseController
  {
    private readonly PolicyService _service;
    private readonly UserService _userService;
    private readonly InodeEngine _inodeEngine;
    private readonly InodeNotificationService _inodeNotificationService;
    private readonly WorkspaceNotificationService _workspaceNotificationService;
    private readonly IMapper _mapper;

    public PoliciesController(
      PolicyService service,
      UserService userService,
      InodeEngine inodeEngine,
      InodeNotificationService inodeNotificationService,
      WorkspaceNotificationService workspaceNotificationService,
      IMapper mapper)
    {
      _service = service;
      _userService = userService;
      _inodeEngine = inodeEngine;
      _inodeNotificationService = inodeNotificationService;
      _workspaceNotificationService = workspaceNotificationService;
      _mapper = mapper;
    }

    [HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
    public async Task<IActionResult> CreateAsync([FromBody] CreatePolicyRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      Policy policy = await _service.InsertAsync(_mapper.Map<Policy>(request), user);
      
      return Ok(_mapper.Map<PolicyDto>(policy));
    }
    
    [HttpDelete("{id}")]
    [SwaggerOperation("Delete", OperationId = "delete")]
    public async Task<IActionResult> DeleteAsync(string id)
    {
      User user = await _userService.FindAsync(UserId);
      
      Policy policy = await _service.FindAsync(id, user);

      await _service.DeleteAsync(policy, user);
      
      return Ok();
    }
    
    [HttpPost("grant")]
    [SwaggerOperation("Grant", OperationId = "grant")]
    public async Task<IActionResult> GrantAsync([FromBody] GrantPolicyRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      await _service.GrantAsync(request.Subject, request.Object, request.Permission, user);

      await DoInodeNotificationsAsync(request.Object, user);

      return Ok();
    }

    [HttpPost("revoke")]
    [SwaggerOperation("Revoke", OperationId = "revoke")]
    public async Task<IActionResult> RevokeAsync([FromBody] RevokePolicyRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      await _service.RevokeAsync(request.Subject, request.Object, request.Permission, user);

      await DoInodeNotificationsAsync(request.Object, user);

      return Ok();
    }
    
    [HttpPost("create")]
    [SwaggerOperation("Create many", OperationId = "createMany")]
    public async Task<IActionResult> CreateManyAsync([FromBody] CreateManyPoliciesRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      foreach (var requestObject in request.Objects)
      {
        try
        {
          await _service.InsertAsync(_mapper.Map<Policy>(requestObject), user);
        }
        catch
        {
          // Ignored
        }
      }

      return Ok();
    }

    [HttpPost("deleteMany")]
    [SwaggerOperation("Delete many", OperationId = "deleteMany")]
    public async Task<IActionResult> DeleteMany([FromBody] DeleteManyPoliciesRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      foreach (var id in request.Ids)
      {
        try
        {
          var role = await _service.FindAsync(id, user);
          await _service.DeleteAsync(role, user);
        }
        catch
        {
          // Ignored
        }
      }

      return Ok();
    }
    
    [HttpGet]
    [SwaggerOperation("Find all", OperationId = "findAll")]
    public async Task<IActionResult> FindAllAsync()
    {
      User user = await _userService.FindAsync(UserId);
      
      IEnumerable<Policy> policies = await _service.FindAllAsync(user);

      return Ok(policies.Select(e => _mapper.Map<PolicyDto>(e)));
    }
    
    [HttpGet("findAllPaged")]
    [SwaggerOperation("Find all paged", OperationId = "findAllPaged")]
    public async Task<IActionResult> FindAllPagedAsync([FromQuery] int page = 1, [FromQuery] int size = 50)
    {
      User user = await _userService.FindAsync(UserId);
      
      PolicyPagedResult policies = await _service.FindAllPagedAsync(user, page, size);

      var dto = _mapper.Map<PolicyPagedResultDto>(policies);

      return Ok(dto);
    }

    [HttpPost("findAllForUser")]
    [SwaggerOperation("Find all for user", OperationId = "findAllForUser")]
    public async Task<IActionResult> FindAllForUserAsync([FromBody] FindAllForUserRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      IEnumerable<Policy> proto = await _service.FindAllForUserAsync(request.Subject, request.Object, user);
      
      return Ok(proto.Select(e => _mapper.Map<PolicyDto>(e)));
    }
    
    [HttpPost("findAllForRole")]
    [SwaggerOperation("Find all for role", OperationId = "findAllForRole")]
    public async Task<IActionResult> FindAllForRoleAsync([FromBody] FindAllForRoleRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      IEnumerable<Policy> proto = await _service.FindAllForRoleAsync(request.RoleId, request.Object, user);
      
      return Ok(proto.Select(e => _mapper.Map<PolicyDto>(e)));
    }

    [HttpPost("getPermissionsForRole")]
    [SwaggerOperation("Get permissions for role", OperationId = "getPermissionsForRole")]
    public async Task<IActionResult> GetPermissionsForRoleAsync([FromBody] GetPermissionsForRoleRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);
      
      IEnumerable<string> proto = await _service.GetObjectPermissionsForRoleAsync(request.RoleId, request.Object, user);
      
      return Ok(proto);
    }
    
    private async Task DoInodeNotificationsAsync(string obj, User user)
    {
      /* In case object is an inode, notify ancestors */
      try
      {
        await _workspaceNotificationService.SendWorkspacesUpdatedAsync();
        
        InodeFacet inode = await _inodeEngine.FindByIdAsync(obj, user);
        IEnumerable<Inode> pathComponents = await _inodeEngine.GetPathAsync(inode, user);
        await _inodeNotificationService.SendInodesChildrenUpdatedAsync(inode.WorkspaceId, pathComponents.Select(e => e.Id));
      }
      catch (Exception)
      {
        // ignored
      }
    }
  }
}