namespace Defyle.WebApi.Inode.Controllers
{
  using System.Collections.Generic;
  using System.IO;
  using Core.Preview.Models;
  using Core.Preview.Services;
  using Microsoft.AspNetCore.Mvc;

  [Route("{path}/tileMaps")]
  public class TileMapsController(TileMapService tileMapService) : Controller
  {
    private readonly TileMapService _tileMapService = tileMapService;

    [HttpGet("zoomLevels")]
    [ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
    public IActionResult GetZoomLevels(string path)
    {
      IEnumerable<ZoomLevel> zoomLevels = _tileMapService.GetZoomLevels(path);
      return Ok(zoomLevels);
    }

    [HttpGet("tiles/{zoomLevel}/{row}/{col}")]
    [ProducesResponseType(typeof(FileResult), 200)]
    public IActionResult DownloadTile(string path, int zoomLevel, int row, int col)
    {
      (Stream stream, string extension) = _tileMapService.GetTileStream(path, zoomLevel, row, col);
      return File(stream, extension);
    }
  }
}