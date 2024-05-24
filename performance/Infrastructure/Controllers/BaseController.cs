namespace Defyle.WebApi.Infrastructure.Controllers
{
  using System.IdentityModel.Tokens.Jwt;
  using System.Security.Claims;
  using Microsoft.AspNetCore.Mvc;
  using Microsoft.AspNetCore.Mvc.ModelBinding;
  using Responses;

  [ProducesResponseType(typeof(ErrorsResponse), 404)]
  [ProducesResponseType(typeof(ErrorsResponse), 403)]
  [ProducesResponseType(typeof(ErrorsResponse), 422)]
  [ProducesResponseType(typeof(ModelStateDictionary), 400)]
	public class BaseController : Controller
	{
    protected string UserId
    {
      get
      {
        var userId = User.FindFirstValue(JwtRegisteredClaimNames.Sub);
        if (!string.IsNullOrWhiteSpace(userId))
        {
         return userId;
        }

        return null;
      }
    }
	}
}