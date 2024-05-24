namespace Defyle.Core.Preview.Services
{
  using System;
  using System.Collections.Generic;
  using System.IO;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Inode.Models;
  using Microsoft.AspNetCore.StaticFiles;
  using Models;
  using Newtonsoft.Json;
  using Newtonsoft.Json.Linq;
  using Storage.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class TileMapService
	{
		private const string MetadataFilename = "meta.json";
		private readonly FileExtensionContentTypeProvider _fileExtensionContentTypeProvider;
		private readonly CoreSettings _coreSettings;
    private readonly FileService _fileService;

    public TileMapService(
			CoreSettings coreSettings,
      FileService fileService)
    {
			_fileExtensionContentTypeProvider = new FileExtensionContentTypeProvider();
			_coreSettings = coreSettings;
      _fileService = fileService;
    }

		public async Task<IEnumerable<ZoomLevel>> GetZoomLevels(Workspace workspace, Inode inode, User user)
		{
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
			{
				throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
			}

			IEnumerable<ZoomLevel> zoomLevels = GetZoomLevels(workspace, file);

			return zoomLevels;
		}

		public async Task<(Stream stream, string extension)> GetTileStreamAsync(Workspace workspace, Inode inode,
      int zoomLevel,int row, int col, User user)
		{
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
			{
				throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
			}

			string path = GetTileImage(workspace, file, zoomLevel, row, col);
			string extension = Path.GetExtension(path);
			string mime = _fileExtensionContentTypeProvider.Mappings[extension];

			if (!System.IO.File.Exists(path))
			{
				throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
			}

			return (new FileStream(path, FileMode.Open), mime);
		}

    public string GetOutputDirectory(Workspace workspace, File file)
		{
			return Path.Combine(
				_coreSettings.DataDirectory,
				workspace.Id,
				"file-data",
				file.Id,
				"tile-map");
		}

		private string GetTileImage(Workspace workspace, File file, int zoomLevel, int row, int col)
		{
			string directory = Path.Combine(GetOutputDirectory(workspace, file), zoomLevel.ToString());
			var files = new DirectoryInfo(directory).GetFiles($"{row}x{col}.*");

			if (files.Length > 0)
			{
				string path = Path.Combine(directory, files[0].Name);

				if (System.IO.File.Exists(path))
				{
					return path;
				}
			}

			throw new Exception("tile file not found");
		}

		private IEnumerable<ZoomLevel> GetZoomLevels(Workspace workspace, File file)
		{
			string globalMetaFile = Path.Combine(GetOutputDirectory(workspace, file), MetadataFilename);

			if (!System.IO.File.Exists(globalMetaFile))
			{
				throw new Exception($"global meta file not found in '{globalMetaFile}'");
			}

			var globalMetaJson = JObject.Parse(System.IO.File.ReadAllText(globalMetaFile));

			var images = new List<ZoomLevel>();
			int zoomLevels = globalMetaJson.Value<int>("zoomLevels");

			for (int i = 0; i < zoomLevels; i++)
			{
				string zoomLevelMetaFile = Path.Combine(
					GetOutputDirectory(workspace, file), i.ToString(),
					MetadataFilename);

				if (!System.IO.File.Exists(zoomLevelMetaFile))
				{
					throw new Exception($"zoom level meta file not found in '{zoomLevelMetaFile}'");
				}

				ZoomLevel zoomLevel = JsonConvert.DeserializeObject<ZoomLevel>(System.IO.File.ReadAllText(zoomLevelMetaFile));

				images.Add(zoomLevel);
			}

			return images;
		}
	}
}