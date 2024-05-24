namespace Defyle.Core.Preview.Services
{
	using System;
	using System.Collections.Generic;
	using System.IO;
	using Microsoft.AspNetCore.StaticFiles;
	using Models;
	using Newtonsoft.Json;
	using Newtonsoft.Json.Linq;

	public class TileService
	{
		private const string MetadataFilename = "meta.json";
		private readonly FileExtensionContentTypeProvider _fileExtensionContentTypeProvider;

		public TileService()
		{
			_fileExtensionContentTypeProvider = new FileExtensionContentTypeProvider();
		}

		public string Create(string path)
		{
			var outputDirectory = GetOutputDirectory(Guid.NewGuid().ToString());
			new TileBuilder(new TileBuilterOptions
			{
				File = path,
				OutputDirectory = outputDirectory,
				Extension = Path.GetExtension(path)
			}).Build();
			return outputDirectory;
		}

		public (Stream stream, string extension) GetTileStream(string path, int zoomLevel, int row, int col)
		{
			string tilePath = GetTileImage(path, zoomLevel, row, col);
			string extension = Path.GetExtension(tilePath);
			string mime = _fileExtensionContentTypeProvider.Mappings[extension];
			if (!File.Exists(tilePath))
			{
				throw new Exception("tile file not found");
			}
			return (new FileStream(tilePath, FileMode.Open), mime);
		}

		private static string GetOutputDirectory(string path)
		{
			return Path.Combine("out", path);
		}

		private static string GetTileImage(string path, int zoomLevel, int row, int col)
		{
			var directory = Path.Combine(GetOutputDirectory(path), zoomLevel.ToString());
			var files = new DirectoryInfo(directory).GetFiles($"{row}x{col}.*");
			if (files.Length > 0)
			{
				string tilePath = Path.Combine(directory, files[0].Name);

				if (File.Exists(tilePath))
				{
					return tilePath;
				}
			}
			throw new Exception("tile file not found");
		}

		public IEnumerable<ZoomLevel> GetZoomLevels(string path)
		{
			var globalMetaFile = Path.Combine(GetOutputDirectory(path), MetadataFilename);
			if (!File.Exists(globalMetaFile))
			{
				throw new Exception($"global meta file not found in '{globalMetaFile}'");
			}

			var globalMetaJson = JObject.Parse(File.ReadAllText(globalMetaFile));
			var images = new List<ZoomLevel>();
			int zoomLevels = globalMetaJson.Value<int>("zoomLevels");

			for (int i = 0; i < zoomLevels; i++)
			{
				var zoomLevelMetaFile = Path.Combine(GetOutputDirectory(path), i.ToString(), MetadataFilename);
				if (!File.Exists(zoomLevelMetaFile))
				{
					throw new Exception($"zoom level meta file not found in '{zoomLevelMetaFile}'");
				}
				ZoomLevel zoomLevel = JsonConvert.DeserializeObject<ZoomLevel>(File.ReadAllText(zoomLevelMetaFile));
				images.Add(zoomLevel);
			}
			return images;
		}
	}
}