namespace Voltaserve.Tiling.Services
{
    using System;
    using System.Collections.Generic;
    using System.IO;
    using Voltaserve.Tiling.Infra;
    using Microsoft.AspNetCore.StaticFiles;
    using Models;
    using Newtonsoft.Json;
    using Newtonsoft.Json.Linq;

    public class ResourceNotFoundException(string message) : Exception(message) { }

    public class TilesService
    {
        private const string MetaFilename = "meta.json";
        private readonly FileExtensionContentTypeProvider _fileExtensionContentTypeProvider;

        public TilesService()
        {
            _fileExtensionContentTypeProvider = new FileExtensionContentTypeProvider();
        }

        public string Create(string path)
        {
            var id = Ids.New();
            new TilesBuilder(new TilesBuilterOptions
            {
                File = path,
                OutputDirectory = GetOutputDirectory(id),
                Extension = Path.GetExtension(path)
            }).Build();
            return id;
        }

        public (Stream stream, string extension) GetTileStream(string path, int zoomLevel, int row, int col)
        {
            string tilePath = GetTileImage(path, zoomLevel, row, col);
            string extension = Path.GetExtension(tilePath);
            string mime = _fileExtensionContentTypeProvider.Mappings[extension];
            return (new FileStream(tilePath, FileMode.Open), mime);
        }

        private static string GetOutputDirectory(string path)
        {
            return Path.Combine("out", path);
        }

        private static string GetTileImage(string path, int zoomLevel, int row, int col)
        {
            var directory = Path.Combine(GetOutputDirectory(path), zoomLevel.ToString());
            if (!Directory.Exists(directory))
            {
                throw new ResourceNotFoundException(directory);
            }
            var files = new DirectoryInfo(directory).GetFiles($"{row}x{col}.*");
            string tilePath = Path.Combine(directory, files[0].Name);
            if (files.Length > 0)
            {
                if (File.Exists(tilePath))
                {
                    return tilePath;
                }
            }
            throw new ResourceNotFoundException(tilePath);
        }

        public IEnumerable<ZoomLevel> GetZoomLevels(string path)
        {
            var metaPath = Path.Combine(GetOutputDirectory(path), MetaFilename);
            if (!File.Exists(metaPath))
            {
                throw new ResourceNotFoundException(metaPath);
            }

            var metaJson = JObject.Parse(File.ReadAllText(metaPath));
            var images = new List<ZoomLevel>();
            int zoomLevels = metaJson.Value<int>("zoomLevels");

            for (int i = 0; i < zoomLevels; i++)
            {
                var zoomLevelPath = Path.Combine(GetOutputDirectory(path), i.ToString(), MetaFilename);
                if (!File.Exists(zoomLevelPath))
                {
                    throw new ResourceNotFoundException(zoomLevelPath);
                }
                var zoomLevel = JsonConvert.DeserializeObject<ZoomLevel>(File.ReadAllText(zoomLevelPath));
                images.Add(zoomLevel);
            }
            return images;
        }
    }
}