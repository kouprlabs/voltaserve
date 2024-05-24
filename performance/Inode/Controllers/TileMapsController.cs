namespace Defyle.WebApi.Inode.Controllers
{
  using System.Collections.Generic;
  using System.IO;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Inode.Models;
  using Core.Inode.Services;
  using Core.Preview.Models;
  using Core.Preview.Services;
  using Core.Workspace.Models;
  using Core.Workspace.Services;
  using Filters;
  using Infrastructure.Controllers;
  using Microsoft.AspNetCore.Authorization;
  using Microsoft.AspNetCore.Mvc;
  using Swashbuckle.AspNetCore.Annotations;

  [Route("workspaces/{workspaceId}/inodes/{id}/tileMaps")]
	[Authorize]
	[PartitionIdCheck]
  [ApiExplorerSettings(GroupName = "Tile maps")]
	public class TileMapsController : BaseController
	{
		private readonly WorkspaceService _workspaceService;
		private readonly InodeService _inodeService;
		private readonly UserService _userService;
		private readonly TileMapService _authorizedTileMapService;
    private readonly IMapper _mapper;

    public TileMapsController(
			WorkspaceService workspaceService,
			InodeService inodeService,
      UserService userService,
			TileMapService authorizedTileMapService,
      IMapper mapper)
		{
			_workspaceService = workspaceService;
			_inodeService = inodeService;
			_userService = userService;
			_authorizedTileMapService = authorizedTileMapService;
      _mapper = mapper;
    }

		[HttpGet("zoomLevels")]
    [SwaggerOperation("Get zoom levels", OperationId = "getZoomLevels")]
		[ProducesResponseType(typeof(IEnumerable<ZoomLevel>), 200)]
		public async Task<IActionResult> GetZoomLevelsAsync(string workspaceId, string id)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      Inode inode = await _inodeService.FindOneAsync(id, user);

      IEnumerable<ZoomLevel> zoomLevels = await _authorizedTileMapService.GetZoomLevels(workspace, inode, user);

      return Ok(zoomLevels);
		}

		[HttpGet("tiles/{zoomLevel}/{row}/{col}")]
    [SwaggerOperation("Download tile", OperationId = "downloadTile")]
		[ProducesResponseType(typeof(FileResult), 200)]
		public async Task<IActionResult> DownloadTileAsync(string workspaceId, string id, int zoomLevel, int row, int col)
		{
      User user = await _userService.FindAsync(UserId);
      Workspace workspace = await _workspaceService.FindAsync(workspaceId, user);
      Inode inode = await _inodeService.FindOneAsync(id, user);

      (Stream stream, string extension) = await _authorizedTileMapService.GetTileStreamAsync(workspace, inode, zoomLevel, row, col, user);

      return File(stream, extension);
		}
	}
}