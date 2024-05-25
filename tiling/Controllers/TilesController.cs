namespace Voltaserve.Tiling.Controllers
{
    using System.Collections.Generic;
    using System.IO;
    using System.Threading.Tasks;
    using Voltaserve.Tiling.Infra;
    using Voltaserve.Tiling.Models;
    using Voltaserve.Tiling.Services;
    using Microsoft.AspNetCore.Http;
    using Microsoft.AspNetCore.Mvc;

    [Route("v2/tiles")]
    public class TilesController(TilesService tilesService) : Controller
    {
        private readonly TilesService _tilesService = tilesService;

        [HttpPost()]
        public async Task<IActionResult> CreateAsync([FromForm] IFormCollection form)
        {
            string path = null;
            try
            {
                var file = form.Files["file"];
                if (file == null || file.Length == 0)
                {
                    return BadRequest("no file uploaded");
                }
                path = Path.Combine(Path.GetTempPath(), Ids.New() + Path.GetExtension(file.FileName));
                using (var stream = new FileStream(path, FileMode.Create))
                {
                    await file.CopyToAsync(stream);
                }
                var metadata = await _tilesService.CreateAsync(path, form["s3_key"], form["s3_bucket"]);
                return Ok(metadata);
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

        [HttpGet("{s3Bucket}/{s3Key}/metadata")]
        [ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
        public async Task<IActionResult> GetMetadataAsync(string s3Bucket, string s3Key)
        {
            try
            {
                Metadata metadata = await _tilesService.GetMetadataAsync(s3Bucket, s3Key);
                return Ok(metadata);
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

        [HttpGet("{s3Bucket}/{s3Key}/zoom_level/{zoomLevel}/row/{row}/col/{col}/ext/{ext}")]
        [ProducesResponseType(typeof(FileResult), 200)]
        public async Task<IActionResult> DownloadTileAsync(string s3Bucket, string s3Key, int zoomLevel, int row, int col, string ext)
        {
            try
            {
                (Stream stream, string contentType) = await _tilesService.GetTileStreamAsync(s3Bucket, s3Key, zoomLevel, row, col, ext);
                return File(stream, contentType);
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