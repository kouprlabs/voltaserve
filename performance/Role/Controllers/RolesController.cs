namespace Defyle.WebApi.Role.Controllers
{
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Infrastructure.Exceptions;
  using Core.Infrastructure.Poco;
  using Core.Role.Models;
  using Core.Role.Pocos;
  using Core.Role.Services;
  using Dtos;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Omu.ValueInjecter;
  using Requests;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("roles")]
  [Authorize]
  [ApiExplorerSettings(GroupName = "Roles")]
  public class RolesController : BaseController
  {
    private readonly RoleService _service;
    private readonly UserService _userService;
    private readonly IMapper _mapper;

    public RolesController(
      RoleService service,
      UserService userService,
      IMapper mapper)
    {
      _service = service;
      _userService = userService;
      _mapper = mapper;
    }

    [HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
    public async Task<IActionResult> CreateAsync([FromBody] CreateRoleRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      Role role = await _service.InsertAsync(_mapper.Map<Role>(request), user);

      return Created("getRoleInformation", _mapper.Map<RoleDto>(role));
    }
    
    [HttpGet("{id}", Name = "getRoleInformation")]
    [SwaggerOperation("Get information", OperationId = "getInformation")]
    public async Task<IActionResult> GetInformationAsync(string id)
    {
      User user = await _userService.FindAsync(UserId);

      var role = await _service.FindAsync(id, user);

      return Ok(_mapper.Map<RoleDto>(role));
    }

    [HttpDelete("{id}")]
    [SwaggerOperation("Delete", OperationId = "delete")]
    public async Task<IActionResult> DeleteAsync(string id)
    {
      User user = await _userService.FindAsync(UserId);

      Role role = await _service.FindAsync(id, user);

      await _service.DeleteAsync(role, user);

      return Ok();
    }
    
    [HttpPatch("{id}")]
    [SwaggerOperation("Update", OperationId = "update")]
    public async Task<IActionResult> UpdateAsync([FromBody] UpdateRoleRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      Role subject = await _service.FindAsync(request.Id, user);
      subject.InjectFrom(request);
      
      Role updated = await _service.UpdateAsync(subject, user);

      return Ok(_mapper.Map<RoleDto>(updated));
    }
    
    [HttpPost("createMany")]
    [SwaggerOperation("Create many", OperationId = "createMany")]
    public async Task<IActionResult> CreateManyAsync([FromBody] CreateManyRolesRequest request)
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
          await _service.InsertAsync(_mapper.Map<Role>(requestObject), user);
        }
        catch
        {
          // Ignored
        }
      }

      return Ok();
    }

    [HttpPost("updateMany")]
    [SwaggerOperation("Update many", OperationId = "updateMany")]
    public async Task<IActionResult> UpdateMany([FromBody] UpdateManyRolesRequest request)
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
          Role subject = await _service.FindAsync(requestObject.Id, user);
          subject.InjectFrom(request);
          
          await _service.UpdateAsync(_mapper.Map<Role>(requestObject), user);
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
    public async Task<IActionResult> DeleteMany([FromBody] DeleteManyRolesRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);

      var errors = new List<Error>();
      
      foreach (var id in request.Ids)
      {
        try
        {
          var role = await _service.FindAsync(id, user);
          await _service.DeleteAsync(role, user);
        }
        catch (GenericException e)
        {
          errors.AddRange(e.Errors);
        }
      }

      if (errors.Any())
      {
        throw new GenericException().WithErrors(errors);
      }

      return Ok();
    }
    
    [HttpGet]
    [SwaggerOperation("Find all", OperationId = "findAll")]
    public async Task<IActionResult> FindAllAsync()
    {
      User user = await _userService.FindAsync(UserId);
      
      var roles = await _service.FindAllAsync(user);

      return Ok(roles.Select(e => _mapper.Map<RoleDto>(e)));
    }

    [HttpGet("findAllPaged")]
    [SwaggerOperation("Find all paged", OperationId = "findAllPaged")]
    public async Task<IActionResult> FindAllPagedAsync([FromQuery] int page = 1, [FromQuery] int size = 50)
    {
      User user = await _userService.FindAsync(UserId);
      
      RolePagedResult roles = await _service.FindAllPagedAsync(user, page, size);

      var dto = _mapper.Map<RolePagedResultDto>(roles);

      return Ok(dto);
    }
    
    [HttpGet("findAllForObject/{obj}")]
    [SwaggerOperation("Find all for object", OperationId = "findAllForObject")]
    public async Task<IActionResult> FindAllForObjectAsync(string obj)
    {
      User user = await _userService.FindAsync(UserId);
      
      IEnumerable<Role> roles = await _service.FindAllForObjectAsync(obj, user);

      return Ok(roles.Select(e => _mapper.Map<RoleDto>(e)));
    }
    
    [HttpGet("findAllWithoutPermissionsForObject/{obj}")]
    [SwaggerOperation("Find all for object", OperationId = "findAllWithoutPermissionsForObject")]
    public async Task<IActionResult> FindAllWithoutPermissionsForObjectAsync(string obj)
    {
      User user = await _userService.FindAsync(UserId);
      
      IEnumerable<Role> roles = await _service.FindAllWithoutPermissionsForObjectAsync(obj, user);

      return Ok(roles.Select(e => _mapper.Map<RoleDto>(e)));
    }
  }
}