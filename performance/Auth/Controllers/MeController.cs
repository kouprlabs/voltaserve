namespace Defyle.WebApi.Auth.Controllers
{
  using System.IO;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Poco;
  using Core.Auth.Services;
  using Core.Infrastructure.Exceptions;
  using Core.Storage.Models;
  using Dtos;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Http;
  using Microsoft.AspNetCore.Mvc;
  using Requests;
  using Responses;
  using Services;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("me")]
  [ApiExplorerSettings(GroupName = "Me")]
  [Authorize]
  public class MeController : BaseController
  {
    private readonly UserService _userService;
    private readonly UserDtoService _userDtoService;
    private readonly IMapper _mapper;
    private readonly PasswordService _passwordService;

    public MeController(
      UserService userService,
      UserDtoService userDtoService,
      IMapper mapper,
      PasswordService passwordService)
    {
      _userService = userService;
      _userDtoService = userDtoService;
      _mapper = mapper;
      _passwordService = passwordService;
    }
    
    [HttpGet]
    [SwaggerOperation("Get information", OperationId = "getInformation")]
    [ProducesResponseType(typeof(UserDto), 200)]
    public async Task<IActionResult> GetInformationAsync()
    {
      try
      {
        User user = await _userService.FindAsync(UserId);
        UserDto dto = await _userDtoService.CreateAsync(user);
        return Ok(dto);
      }
      catch (ResourceNotFoundException)
      {
        return NotFound();
      }
    }
    
    [HttpPatch]
    [SwaggerOperation("Update", OperationId = "update")]
    [ProducesResponseType(typeof(UpdateUserResponse), 200)]
    [ProducesResponseType(typeof(UpdateUserResponse), 207)]
    public async Task<IActionResult> UpdateAsync([FromBody] UpdateAccountRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _userService.FindAsync(UserId);
      
      User subject = await _userService.FindAsync(UserId);
      
      // TODO: move this security check to the user service
      User newSubject = await _userService.FindAsync(UserId);
      if (!string.IsNullOrEmpty(request.CurrentPassword) && !string.IsNullOrEmpty(request.NewPassword))
      {
        if (!_passwordService.VerifyHashedPassword(user.PasswordHash, request.CurrentPassword))
        {
          return Forbid();
        }
        newSubject.PasswordHash = request.NewPassword;
      }
      if (!string.IsNullOrWhiteSpace(request.NewEmail))
      {
        newSubject.Email = request.NewEmail.ToLowerInvariant().Trim();
      }
      if (!string.IsNullOrWhiteSpace(request.NewFullName))
      {
        newSubject.FullName = request.NewFullName.Trim();
      }
      
      UpdateUserResult result = await _userService.UpdateAsync(subject, newSubject, user);

      var response = _mapper.Map<UpdateUserResponse>(result);
      if (response.Errors.Count > 0)
      {
        return new ObjectResult(response) {StatusCode = 207};
      }
      else
      {
        return Ok(response);
      }
    }
    
    [HttpPost("updateImage")]
    [SwaggerOperation("Update image", OperationId = "updateImage")]
    [ProducesResponseType(typeof(UserDto), 200)]
    public async Task<IActionResult> UpdateImageAsync(IFormFile file)
    {
      try
      {
        string path = Path.GetTempFileName();

        using (var stream = new FileStream(path, FileMode.Create))
        {
          file.CopyTo(stream);
        }

        User user = await _userService.FindAsync(UserId);
        User updatedUser = await _userService.UpdateImageAsync(path, user, user);

        return Ok(_mapper.Map<UserDto>(updatedUser));
      }
      catch (ResourceNotFoundException)
      {
        return NotFound();
      }
    }

    [HttpPost("updateImageFromGravatar")]
    [SwaggerOperation("Update image from Gravatar", OperationId = "updateImageFromGravatar")]
    [ProducesResponseType(typeof(UserDto), 200)]
    public async Task<IActionResult> UpdateImageFromGravatarAsync()
    {
      User user = await _userService.FindAsync(UserId);
      User updatedUser = await _userService.RefreshImageFromGravatarAsync(user, user);
      return Ok(_mapper.Map<UserDto>(updatedUser));
    }

    [HttpDelete]
    [SwaggerOperation("Delete", OperationId = "delete")]
    [ProducesResponseType(typeof(void), 200)]
    public async Task<IActionResult> DeleteAsync()
    {
      User user = await _userService.FindAsync(UserId);
      await _userService.DeleteAsync(user, user);
      return Ok();
    }
    
    [HttpGet("getStorageUsage")]
    [SwaggerOperation("Get storage usage", OperationId = "getStorageUsage")]
    [ProducesResponseType(typeof(StorageUsage), 200)]
    public async Task<IActionResult> GetStorageUsageAsync()
    {
      User user = await _userService.FindAsync(UserId);
      StorageUsage storageUsage = await _userService.GetStorageUsageAsync(user);
      return Ok(storageUsage);
    }
  }
}