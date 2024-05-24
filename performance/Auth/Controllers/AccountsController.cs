namespace Defyle.WebApi.Auth.Controllers
{
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Dtos;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Requests;
  using Responses;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("accounts")]
	[Authorize]
  [ApiExplorerSettings(GroupName = "Accounts")]
	public class AccountsController : BaseController
	{
		private readonly AccountService _accountService;
    private readonly IMapper _mapper;

    public AccountsController(
      AccountService accountService,
      IMapper mapper)
    {
      _accountService = accountService;
      _mapper = mapper;
    }

		
		[HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
    [ProducesResponseType(typeof(UserDto), 200)]
    [AllowAnonymous]
		public async Task<IActionResult> CreateAsync([FromBody] CreateAccountRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _accountService.CreateLocalAsync(_mapper.Map<User>(request));
      return Ok(_mapper.Map<UserDto>(user));
		}

    [HttpPost("resetPassword")]
    [SwaggerOperation("Reset password", OperationId = "resetPassword")]
		[AllowAnonymous]
		public async Task<IActionResult> ResetPasswordAsync([FromBody] ResetPasswordRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      await _accountService.ResetPasswordAsync(request.Token, request.NewPassword);
      return Ok();
		}

		[HttpPost("confirmEmail")]
    [SwaggerOperation("Confirm email", OperationId = "confirmEmail")]
		[AllowAnonymous]
		public async Task<IActionResult> ConfirmEmailAsync([FromBody] ConfirmEmailRequest request)
		{
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      User user = await _accountService.ConfirmEmailAsync(request.Token);
      return Ok(new EmailConfirmationResponse{Username = user.Username});
		}
    
    [HttpPost("sendResetPasswordEmail")]
    [SwaggerOperation("Send reset password email", OperationId = "sendResetPasswordEmail")]
    [AllowAnonymous]
    public async Task<IActionResult> SendResetPasswordEmailAsync([FromBody] SendResetPasswordEmailRequest request)
    {
      if (!ModelState.IsValid)
      {
        return BadRequest(ModelState);
      }
      
      await _accountService.SendResetPasswordEmailAsync(request.Email);
      return Ok();
    }
	}
}