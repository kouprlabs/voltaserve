namespace Defyle.Core.Infrastructure.Services
{
  using System;
  using Exceptions;
  using Grpc.Core;
  using Microsoft.AspNetCore.Hosting;
  using Microsoft.Extensions.Hosting;
  using Newtonsoft.Json;
  using Poco;

  public class GrpcExceptionTranslator
  {
    private readonly IWebHostEnvironment _env;

    public GrpcExceptionTranslator(IWebHostEnvironment env)
    {
      _env = env;
    }
    
    public Exception Translate(Exception exception)
    {
      if (exception is RpcException rpcException)
      {
        var jsonError = JsonConvert.DeserializeObject<JsonError>(rpcException.Status.Detail);
        
        string description = _env.IsDevelopment() ? jsonError.InternalDescription : jsonError.PublicDescription;

        Error error = new Error()
          .WithCode(jsonError.Code)
          .WithDescription(description);

        switch (jsonError.Code)
        {
          case "resource_not_found":
            return new ResourceNotFoundException().WithError(error);

          case "internal_server_error":
            return new InternalServerErrorException().WithError(error);;

          default:
            return new GenericException().WithError(error);
        }
      }
      else
      {
        return exception;
      }
    }
  }
}