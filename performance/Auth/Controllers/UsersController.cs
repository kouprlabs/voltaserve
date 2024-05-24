namespace Defyle.WebApi.Auth.Controllers
{
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Poco;
  using Core.Auth.Services;
  using Dtos;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Omu.ValueInjecter;
  using Requests;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("users")]
  [Authorize]
  [ApiExplorerSettings(GroupName = "Users")]
  public class UsersController : BaseController
  {
    private readonly UserService _service;
    private readonly IMapper _mapper;

    public UsersController(
      UserService service,
      IMapper mapper)
    {
      _service = service;
      _mapper = mapper;
    }
    
    [HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
    [ProducesResponseType(typeof(UserDto), 200)]
    public async Task<IActionResult> CreateAsync([FromBody] CreateUserRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _service.FindAsync(UserId);
        
      User newUser = await _service.InsertAsync(_mapper.Map<User>(request), user);
        
      return Ok(_mapper.Map<UserDto>(newUser));
    }

    [HttpGet("{id}")]
    [SwaggerOperation("Get information", OperationId = "getInformation")]
    [ProducesResponseType(typeof(UserDto), 200)]
    public async Task<IActionResult> GetInformation(string id)
    {
      User user = await _service.FindAsync(UserId);
        
      User newUser = await _service.FindAsync(id, user);
        
      return Ok(_mapper.Map<UserDto>(newUser));
    }
    
    [HttpDelete("{id}")]
    [SwaggerOperation("Delete", OperationId = "delete")]
    [ProducesResponseType(typeof(void), 200)]
    public async Task<IActionResult> DeleteAsync(string id)
    {
      User user = await _service.FindAsync(id);
      
      User subject = await _service.FindAsync(UserId);
      
      await _service.DeleteAsync(subject, user);
      
      return Ok();
    }
    
    [HttpPatch("{id}")]
    [SwaggerOperation("Update", OperationId = "update")]
    public async Task<IActionResult> UpdateAsync([FromBody] UpdateUserRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _service.FindAsync(UserId);

      User subject = await _service.FindAsync(request.Id);
      
      User newSubject = await _service.FindAsync(request.Id);
      newSubject.InjectFrom(request);
      
      UpdateUserResult result = await _service.UpdateAsync(subject, newSubject, user);

      return Ok(_mapper.Map<UserDto>(result.Result));
    }
    
    [HttpPost("create")]
    [SwaggerOperation("Create many", OperationId = "createMany")]
    public async Task<IActionResult> CreateManyAsync([FromBody] CreateManyUsersRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _service.FindAsync(UserId);

      foreach (var requestObject in request.Objects)
      {
        try
        {
          await _service.InsertAsync(_mapper.Map<User>(requestObject), user);
        }
        catch
        {
          // Ignored
        }
      }

      return Ok();
    }

    [HttpPost("update")]
    [SwaggerOperation("Update many", OperationId = "updateMany")]
    public async Task<IActionResult> UpdateMany([FromBody] UpdateManyUsersRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _service.FindAsync(UserId);

      foreach (var requestObject in request.Objects)
      {
        try
        {
          User subject = await _service.FindAsync(requestObject.Id);
          
          User newSubject = await _service.FindAsync(requestObject.Id);
          newSubject.InjectFrom(request);
          
          await _service.UpdateAsync(subject, newSubject, user);
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
    public async Task<IActionResult> DeleteMany([FromBody] DeleteManyUsersRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _service.FindAsync(UserId);

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
    [ProducesResponseType(typeof(IEnumerable<UserDto>), 200)]
    public async Task<IActionResult> FindAllAsync()
    {
      User user = await _service.FindAsync(UserId);
        
      IEnumerable<User> users = await _service.FindAllAsync(user);

      return Ok(users.Select(e => _mapper.Map<UserDto>(e)));
    }
    
    [HttpGet("findAllPaged")]
    [SwaggerOperation("Find all paged", OperationId = "findAllPaged")]
    public async Task<IActionResult> FindAllPagedAsync([FromQuery] int page = 1, [FromQuery] int size = 50)
    {
      User user = await _service.FindAsync(UserId);
      
      var roles = await _service.FindAllPagedAsync(user, page, size);

      var dto = _mapper.Map<UserPagedResultDto>(roles);

      return Ok(dto);
    }
  }
}