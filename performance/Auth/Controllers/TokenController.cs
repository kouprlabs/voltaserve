namespace Defyle.WebApi.Auth.Controllers
{
  using System;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Exceptions;
  using Core.Auth.Models;
  using Core.Auth.Poco;
  using Core.Auth.Services;
  using Dtos;
  using Infrastructure.Dtos;
  using Infrastructure.Responses;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Requests;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("token")]
  [ApiExplorerSettings(GroupName = "Token")]
	public class TokenController : Controller
	{
		private readonly TokenService _tokenService;
    private readonly IMapper _mapper;

    public TokenController(
			TokenService tokenService,
      IMapper mapper)
    {
      _tokenService = tokenService;
      _mapper = mapper;
    }

		[AllowAnonymous]
		[HttpPost]
    [SwaggerOperation("Create", OperationId = "create")]
		[Consumes("application/x-www-form-urlencoded")]
		[Produces("application/json")]
		[ProducesResponseType(typeof(TokenDto), 200)]
		[ProducesResponseType(typeof(ErrorsResponse), 401)]
		public async Task<IActionResult> ExchangeAsync([FromForm] TokenExchangeRequest request)
		{
			try
			{
				Token token = await _tokenService.ExchangeAsync(_mapper.Map<TokenExchangeOptions>(request));
				return Ok(_mapper.Map<Token, TokenDto>(token));
			}
			catch (AuthenticationException e)
			{
				return new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 401};
			}
			catch (EmailNotConfirmedException e)
			{
				return new ObjectResult(new ErrorsResponse{Errors = e.Errors.Select(x => _mapper.Map<ErrorDto>(x))}) {StatusCode = 403};
			}
			catch (ArgumentException e)
			{
				return BadRequest(e.Message);
			}
			catch (Exception)
			{
				return Unauthorized();
			}
		}
	}
}