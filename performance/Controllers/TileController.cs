namespace Voltaserve.Performance.Controllers
{
    using System;
    using System.Collections.Generic;
    using System.IO;
    using System.Threading.Tasks;
    using Defyle.Performance.Infra;
    using Defyle.Performance.Models;
    using Defyle.Performance.Services;
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
                    return BadRequest("no file uploaded");
                }
                path = Path.Combine(Path.GetTempPath(), Ids.New() + Path.GetExtension(file.FileName));
                using (var stream = new FileStream(path, FileMode.Create))
                {
                    await file.CopyToAsync(stream);
                }
                var id = _tileService.Create(path);
                return Ok(id);
            }
            catch
            {
                return StatusCode(500);
            }
            finally
            {
                if (path != null && System.IO.File.Exists(path))
                {
                    System.IO.File.Delete(path);
                }
            }
        }

        [HttpGet("{path}/zoom_levels")]
        [ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
        public IActionResult GetZoomLevels(string path)
        {
            try
            {
                IEnumerable<ZoomLevel> zoomLevels = _tileService.GetZoomLevels(path);
                return Ok(zoomLevels);
            }
            catch (ResourceNotFoundException)
            {
                return NotFound();
            }
            catch
            {
                return StatusCode(500);
            }
        }

        [HttpGet("{path}/zoom_level/{zoomLevel}/row/{row}/col/{col}")]
        [ProducesResponseType(typeof(FileResult), 200)]
        public IActionResult DownloadTile(string path, int zoomLevel, int row, int col)
        {
            try
            {
                (Stream stream, string extension) = _tileService.GetTileStream(path, zoomLevel, row, col);
                return File(stream, extension);
            }
            catch (ResourceNotFoundException)
            {
                return NotFound();
            }
            catch
            {
                return StatusCode(500);
            }
        }
    }
}