namespace Defyle.WebApi.Configuration
{
  using System.IdentityModel.Tokens.Jwt;
  using System.Threading.Tasks;
  using Core.Auth.Services;
  using Microsoft.AspNetCore.Authentication.JwtBearer;
  using Microsoft.Extensions.DependencyInjection;
  using Microsoft.IdentityModel.Tokens;

  public static class SecurityExtensions
  {
    public static void AddSecurity(this IServiceCollection services)
    {
      JwtSecurityTokenHandler.DefaultInboundClaimTypeMap.Clear();

      services.AddAuthentication(options =>
        {
          options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
          options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
          options.DefaultSignInScheme = JwtBearerDefaults.AuthenticationScheme;
        })
        .AddJwtBearer(options =>
        {
          TokenService tokenService = services.BuildServiceProvider().GetService<TokenService>();
          TokenValidationParameters tokenValidationParameters = tokenService.TokenValidationParameters;

          options.TokenValidationParameters = tokenValidationParameters;

          // We have to hook the OnMessageReceived event in order to
          // allow the JWT authentication handler to read the access
          // token from the query string when a WebSocket or
          // Server-Sent Events request comes in.
          options.Events = new JwtBearerEvents
          {
            OnMessageReceived = context =>
            {
              var accessToken = context.Request.Query["access_token"];

              // If the request is for our hub...
              var path = context.HttpContext.Request.Path;
              if (!string.IsNullOrEmpty(accessToken) && path.StartsWithSegments("/notifications"))
              {
                // Read the token out of the query string
                context.Token = accessToken;
              }

              return Task.CompletedTask;
            }
          };
        });
    }
  }
}