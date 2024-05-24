namespace Defyle.WebApi.Configuration
{
  using System;
  using System.Linq;
  using AutoMapper;
  using Core.Auth.Exceptions;
  using Core.Infrastructure.Exceptions;
  using Core.Infrastructure.Poco;
  using Core.Storage.Exceptions;
  using Grpc.Core;
  using Infrastructure.Dtos;
  using Infrastructure.Responses;
  using Microsoft.AspNetCore.Builder;
  using Microsoft.AspNetCore.Diagnostics;
  using Microsoft.AspNetCore.Hosting;
  using Microsoft.AspNetCore.Http;
  using Microsoft.Extensions.Hosting;
  using Newtonsoft.Json;

  public static class ExceptionHandlerExtensions
  {
    public static void UseExceptionHandler(this IApplicationBuilder app, IWebHostEnvironment env, IMapper mapper)
    {
      app.UseExceptionHandler(errorApp =>
      {
        errorApp.Run(async context =>
        {
          context.Response.ContentType = "application/json";
          
          var exception = context.Features.Get<IExceptionHandlerPathFeature>()?.Error;
          if (exception is RpcException rpcException)
          {
            try
            {
              var jsonError = JsonConvert.DeserializeObject<JsonError>(rpcException.Status.Detail);
              switch (jsonError.Code)
              {
                case "resource_not_found":
                  context.Response.StatusCode = 404;
                  break;
                
                case "internal_server_error":
                  context.Response.StatusCode = 500;
                  break;
                
                default:
                  context.Response.StatusCode = 500;
                  break;
              }

              string description = env.IsDevelopment() ? jsonError.InternalDescription : jsonError.PublicDescription;

              var errorsResponse = new ErrorsResponse
              {
                Errors = new[] {new ErrorDto {Code = jsonError.Code, Description = description}}
              };

              await context.Response.WriteAsync(JsonConvert.SerializeObject(errorsResponse));
            }
            catch
            {
              context.Response.StatusCode = 500;
            }
          }
          else if (exception is GenericException baseException)
          {
            var errorsResponse = new ErrorsResponse {Errors = baseException.Errors.Select(mapper.Map<ErrorDto>)};
            
            if (exception is ResourceNotFoundException)
            {
              context.Response.StatusCode = 404;
            }
            else if (exception is InternalServerErrorException)
            {
              context.Response.StatusCode = 500;
            }
            else if (exception is AuthenticationException)
            {
              context.Response.StatusCode = 401;
            }
            else if (exception is EmailNotConfirmedException)
            {
              context.Response.StatusCode = 403;
            }
            else if (exception is EmailValidationException)
            {
              context.Response.StatusCode = 422;
            }
            else if (exception is PasswordValidationException)
            {
              context.Response.StatusCode = 422;
            }
            else if (exception is StorageUsageExceededException)
            {
              context.Response.StatusCode = 403;
            }
            else
            {
              context.Response.StatusCode = 500;
            }
            
            await context.Response.WriteAsync(JsonConvert.SerializeObject(errorsResponse));
          }
          else if (exception is ArgumentException)
          {
            context.Response.StatusCode = 400;
          }
          else
          {
            context.Response.StatusCode = 500;
          }
        });
      });
    }
  }
}