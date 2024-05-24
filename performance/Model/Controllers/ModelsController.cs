namespace Defyle.WebApi.Model.Controllers
{
  using System.Threading.Tasks;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Model.Services;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Requests;

  [Route("models")]
  [Authorize]
  [ApiExplorerSettings(GroupName = "Models")]
  public class ModelsController : BaseController
  {
    private readonly ModelService _service;
    private readonly UserService _userService;

    public ModelsController(
      ModelService service,
      UserService userService)
    {
      _service = service;
      _userService = userService;
    }
    
    [HttpGet("getPermissions/{modelName}")]
    public async Task<IActionResult> GetPermissions(string modelName)
    {
      return Ok(await _service.GetPermissionsAsync(modelName));
    }
    
    [HttpPost("getProperty")]
    [Produces("application/json")]
    public async Task<IActionResult> GetColumnDefs([FromBody] GetPropertyRequest request)
    {
      User user = await _userService.FindAsync(UserId);

      return Ok(await _service.GetPropertyAsync(request.Model, request.Property, user));
    }
  }
}