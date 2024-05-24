namespace Defyle.WebApi.Filters
{
  using AutoMapper;
  using Core.Infrastructure.Poco;
  using Infrastructure.Dtos;
  using Infrastructure.Responses;
  using Microsoft.AspNetCore.Mvc;
	using Microsoft.AspNetCore.Mvc.Filters;

  public class PartitionIdCheckAttribute : TypeFilterAttribute
	{
    public PartitionIdCheckAttribute() : base(typeof(PartitionIdCheckAttributeImpl))
    {
    }

		private class PartitionIdCheckAttributeImpl : IActionFilter
		{
      private readonly IMapper _mapper;

      public PartitionIdCheckAttributeImpl(IMapper mapper)
      {
        _mapper = mapper;
      }
      
			public void OnActionExecuting(ActionExecutingContext context)
			{
				var partitionId = context.HttpContext.Request.Headers["PartitionId"];

				if (string.IsNullOrWhiteSpace(partitionId))
				{
					context.Result = new BadRequestObjectResult(
            new ErrorsResponse{Errors = new []{_mapper.Map<ErrorDto>(Error.MissingPartitionIdHeaderError)}});
				}
			}

			public void OnActionExecuted(ActionExecutedContext context)
			{
			}
		}
	}
}