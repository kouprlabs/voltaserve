namespace Voltaserve.Mosaic.Controllers
{
    using System.Collections.Generic;
    using System.IO;
    using System.Threading.Tasks;
    using Voltaserve.Mosaic.Infra;
    using Voltaserve.Mosaic.Models;
    using Voltaserve.Mosaic.Services;
    using Microsoft.AspNetCore.Http;
    using Microsoft.AspNetCore.Mvc;
    using Microsoft.Extensions.Logging;
    using System;

    [Route("v2/mosaics")]
    public class MosaicController(MosaicService _mosaicService, ILogger<MosaicController> _logger) : Controller
    {
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
                var metadata = await _mosaicService.CreateAsync(path, form["s3_key"], form["s3_bucket"]);
                return Ok(metadata);
            }
            catch (Exception e)
            {
                _logger.LogError(e, "Failed to create mosaic.");
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

        [HttpDelete("{s3Bucket}/{s3Key}")]
        public async Task<IActionResult> DeleteAsync(string s3Bucket, string s3Key)
        {
            try
            {
                await _mosaicService.DeleteAsync(s3Bucket, s3Key);
                return NoContent();
            }
            catch (ResourceNotFoundException)
            {
                return NotFound();
            }
            catch (Exception e)
            {
                _logger.LogError(e, "Failed to delete mosaic.");
                return StatusCode(500);
            }
        }

        [HttpGet("{s3Bucket}/{s3Key}/metadata")]
        [ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
        public async Task<IActionResult> GetMetadataAsync(string s3Bucket, string s3Key)
        {
            try
            {
                Metadata metadata = await _mosaicService.GetMetadataAsync(s3Bucket, s3Key);
                return Ok(metadata);
            }
            catch (ResourceNotFoundException)
            {
                return NotFound();
            }
            catch (Exception e)
            {
                _logger.LogError(e, "Failed to get mosaic metadata.");
                return StatusCode(500);
            }
        }

        [HttpGet("{s3Bucket}/{s3Key}/zoom_level/{zoomLevel}/row/{row}/col/{col}/ext/{ext}")]
        [ProducesResponseType(typeof(FileResult), 200)]
        public async Task<IActionResult> DownloadTileAsync(string s3Bucket, string s3Key, int zoomLevel, int row, int col, string ext)
        {
            try
            {
                (Stream stream, string contentType) = await _mosaicService.GetTileStreamAsync(s3Bucket, s3Key, zoomLevel, row, col, ext);
                return File(stream, contentType);
            }
            catch (ResourceNotFoundException)
            {
                return NotFound();
            }
            catch (Exception e)
            {
                _logger.LogError(e, "Failed to get mosaic.");
                return StatusCode(500);
            }
        }
    }
}