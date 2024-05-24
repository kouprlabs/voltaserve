namespace Defyle.WebApi.Inode.Controllers
{
  using System;
  using System.Collections.Generic;
  using System.IO;
  using System.Threading.Tasks;
  using Core.Preview.Models;
  using Core.Preview.Services;
  using Microsoft.AspNetCore.Http;
  using Microsoft.AspNetCore.Mvc;

  [Route("v2/performance/tiles")]
  public class TileController(TileService tileService) : Controller
  {
    private readonly TileService _tileService = tileService;

    [HttpPost()]
    public async Task<IActionResult> Create(IFormFile file)
    {
      string path = null;
      try
      {
        if (file == null || file.Length == 0)
        {
          return BadRequest("No file uploaded");
        }
        path = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString() + Path.GetExtension(file.FileName));
        using (var stream = new FileStream(path, FileMode.Create))
        {
          await file.CopyToAsync(stream);
        }
        var outputDirectory = _tileService.Create(path);
        return Ok(outputDirectory);
      }
      catch
      {
        return StatusCode(500);
      }
      finally
      {
        if (path != null)
        {
          if (System.IO.File.Exists(path))
          {
            System.IO.File.Delete(path);
          }
        }
      }
    }

    [HttpGet("{path}/zoom_levels")]
    [ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
    public IActionResult GetZoomLevels(string path)
    {
      IEnumerable<ZoomLevel> zoomLevels = _tileService.GetZoomLevels(path);
      return Ok(zoomLevels);
    }

    [HttpGet("{path}/{zoomLevel}/{row}/{col}")]
    [ProducesResponseType(typeof(FileResult), 200)]
    public IActionResult DownloadTile(string path, int zoomLevel, int row, int col)
    {
      (Stream stream, string extension) = _tileService.GetTileStream(path, zoomLevel, row, col);
      return File(stream, extension);
    }
  }
}