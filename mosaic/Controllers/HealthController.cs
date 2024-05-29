namespace Voltaserve.Mosaic.Controllers
{
  using Microsoft.AspNetCore.Mvc;

  [Route("v2/health")]
  public class HealthController : Controller
  {
    [HttpGet]
    public IActionResult Health()
    {
      return Ok();
    }
  }
}