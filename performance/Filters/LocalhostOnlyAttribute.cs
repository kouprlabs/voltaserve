namespace Defyle.WebApi.Filters
{
	using Microsoft.AspNetCore.Mvc;
	using Microsoft.AspNetCore.Mvc.Filters;

	public class LocalhostOnlyAttribute : TypeFilterAttribute
	{
		public LocalhostOnlyAttribute() : base(typeof(LocalhostOnlyAttributeImpl))
		{
		}

		private class LocalhostOnlyAttributeImpl : IActionFilter
		{
			public void OnActionExecuting(ActionExecutingContext context)
			{
				if (context.HttpContext.Request.Host.Host != "localhost")
				{
					context.Result = new ForbidResult();
				}
			}

			public void OnActionExecuted(ActionExecutedContext context)
			{
			}
		}
	}
}