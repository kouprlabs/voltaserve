namespace Voltaserve.Tiling.Services
{
    using System;
    using System.Collections.Generic;
    using System.Drawing;
    using System.IO;
    using System.Text.Json;
    using Models;

    public enum ActionOnExistingDirectory
    {
        Delete,
        Skip
    }

    public class TilesBuilterOptions
    {
        public string File { get; set; }

        public string OutputDirectory { get; set; }
    }

    public class TilesBuilder(TilesBuilterOptions options)
    {
        private IImage _image;

        private ScaleDownPercentage _scaleDownPercentage;

        private MinimumScaleSize _minimumScaleSize;

        private TileSize _tileSize;

        private readonly TilesBuilterOptions _options = options;

        public ScaleDownPercentage ScaleDownPercentage
        {
            get
            {
                _scaleDownPercentage ??= new ScaleDownPercentage(70);
                return _scaleDownPercentage;
            }

            set => _scaleDownPercentage = value;
        }

        public MinimumScaleSize MinimumScaleSize
        {
            get
            {
                _minimumScaleSize ??= new MinimumScaleSize(new Size(500, 500));
                return _minimumScaleSize;
            }

            set => _minimumScaleSize = value;
        }

        public TileSize TileSize
        {
            get
            {
                _tileSize ??= new TileSize(new Size(300, 300));
                return _tileSize;
            }

            set => _tileSize = value;
        }

        public ActionOnExistingDirectory ActionOnExistingDirectory { get; set; }

        public Metadata Build()
        {
            bool cleanupIfFails = false;
            if (!Directory.Exists(_options.OutputDirectory))
            {
                Directory.CreateDirectory(_options.OutputDirectory);
                cleanupIfFails = true;
            }
            try
            {
                _image = new Image(_options.File);
                var zoomLevelsIndexes = RequiredZoomLevelIndexes();
                if (zoomLevelsIndexes.Count == 0)
                {
                    throw new Exception("creating zoom levels is not required for this image.");
                }
                var zoomLevels = new List<ZoomLevel>();
                foreach (int index in zoomLevelsIndexes)
                {
                    CreateZoomLevelDirectory(index);
                    Image scaled = Scale(index);
                    var zoomLevel = Decompose(scaled, index, new Region());
                    zoomLevels.Add(zoomLevel);
                }
                var metadata = new Metadata
                {
                    Width = _image.Width,
                    Height = _image.Height,
                    Extension = Path.GetExtension(_options.File),
                    ZoomLevels = zoomLevels,
                };
                File.WriteAllText(
                    GetMetadataFilePath(),
                    JsonSerializer.Serialize(metadata, new JsonSerializerOptions() { WriteIndented = true }));
                return metadata;
            }
            catch
            {
                if (cleanupIfFails)
                {
                    Directory.Delete(_options.OutputDirectory, true);
                }
                throw;
            }
        }

        private static void DeleteDirectoryWithContent(string directory)
        {
            string[] files = Directory.GetFiles(directory);
            string[] dirs = Directory.GetDirectories(directory);
            foreach (string file in files)
            {
                File.SetAttributes(file, FileAttributes.Normal);
                File.Delete(file);
            }
            foreach (string dir in dirs)
            {
                DeleteDirectoryWithContent(dir);
            }
            Directory.Delete(directory, true);
        }

        private ZoomLevel Decompose(Image image, int zoomLevel, Region region)
        {
            bool tileWidthExceeded = image.Width > TileSize.Width;
            bool tileHeightExceeded = image.Height > TileSize.Height;

            int cols = tileWidthExceeded ? image.Width / TileSize.Width : 1;
            int rows = tileHeightExceeded ? image.Height / TileSize.Height : 1;
            int remainingWidth = tileWidthExceeded ? image.Width - (TileSize.Width * cols) : 0;
            int remainingHeight = tileHeightExceeded ? image.Height - (TileSize.Height * rows) : 0;
            int totalCols = remainingWidth != 0 ? cols + 1 : cols;
            int totalRows = remainingHeight != 0 ? rows + 1 : rows;

            TileSize adaptedTileSize = TileSize;
            if (!tileWidthExceeded)
            {
                adaptedTileSize.Width = image.Width;
            }
            if (!tileHeightExceeded)
            {
                adaptedTileSize.Height = image.Height;
            }

            int colStart, colEnd, rowStart, rowEnd;
            bool includesRemainingTiles;

            if (region.IsNull())
            {
                colStart = 0;
                colEnd = cols - 1;
                rowStart = 0;
                rowEnd = rows - 1;
                includesRemainingTiles = true;
            }
            else
            {
                colStart = region.ColStart;
                colEnd = region.ColEnd;
                rowStart = region.RowStart;
                rowEnd = region.RowEnd;
                includesRemainingTiles = region.IncludesRemainingTiles;
            }

            for (int c = colStart; c <= colEnd; c++)
            {
                for (int r = rowStart; r <= rowEnd; r++)
                {
                    var tileMetadata = new TileMetadata
                    {
                        X = c * _tileSize.Width,
                        Y = r * _tileSize.Height,
                        Width = _tileSize.Width,
                        Height = _tileSize.Height,
                        Row = r,
                        Col = c
                    };
                    var clippingRect = new Rectangle
                    {
                        X = tileMetadata.X,
                        Y = tileMetadata.Y,
                        Width = tileMetadata.Width,
                        Height = tileMetadata.Height
                    };
                    var cropped = new Image(image);
                    cropped.Crop(clippingRect);
                    cropped.Save(GetTileOutputPath(zoomLevel, tileMetadata.Row, tileMetadata.Col));
                }
            }

            /* Remaining height */
            if (includesRemainingTiles && remainingHeight > 0)
            {
                for (int c = 0; c < cols; c++)
                {
                    var clippingRect = new Rectangle
                    {
                        X = c * _tileSize.Width,
                        Y = image.Height - remainingHeight,
                        Width = _tileSize.Width,
                        Height = remainingHeight
                    };

                    var cropped = new Image(image);
                    cropped.Crop(clippingRect);
                    cropped.Save(GetTileOutputPath(zoomLevel, totalRows - 1, c));
                }
            }

            /* Remaining width */
            if (includesRemainingTiles && remainingWidth > 0)
            {
                for (int r = 0; r < rows; r++)
                {
                    var clippingRect = new Rectangle
                    {
                        X = image.Width - remainingWidth,
                        Y = r * _tileSize.Height,
                        Width = remainingWidth,
                        Height = _tileSize.Height
                    };

                    var cropped = new Image(image);
                    cropped.Crop(clippingRect);
                    cropped.Save(GetTileOutputPath(zoomLevel, r, totalCols - 1));
                }
            }

            /* Remaining bottom right corner */
            if (includesRemainingTiles && remainingWidth > 0 && remainingHeight > 0)
            {
                var clippingRect = new Rectangle
                {
                    X = image.Width - remainingWidth,
                    Y = image.Height - remainingHeight,
                    Width = remainingWidth,
                    Height = remainingHeight
                };

                var cropped = new Image(image);
                cropped.Crop(clippingRect);
                cropped.Save(GetTileOutputPath(zoomLevel, totalRows - 1, totalCols - 1));
            }

            return new ZoomLevel
            {
                Index = zoomLevel,
                Width = image.Width,
                Height = image.Height,
                Rows = totalRows,
                Cols = totalCols,
                ScaleDownPercentage = GetScaleDownPercentage(zoomLevel),
                Tile = new Tile
                {
                    Width = adaptedTileSize.Width,
                    Height = adaptedTileSize.Height,
                    LastColWidth = remainingWidth,
                    LastRowHeight = remainingHeight
                }
            };
        }

        private static float GetScaleDownPercentage(int zoomLevel)
        {
            float value = 100.0f;
            for (int i = 0; i < zoomLevel; i++)
            {
                value *= 0.70f;
            }
            return value;
        }

        private Image Scale(int zoomLevel)
        {
            Size imageSizeForZoomLevel = GetImageSizeForZoomLevel(zoomLevel);
            var scaled = new Image(_image);
            scaled.ScaleWithAspectRatio(imageSizeForZoomLevel.Width, imageSizeForZoomLevel.Height);
            return scaled;
        }

        private Size GetImageSizeForZoomLevel(int zoomLevel)
        {
            var size = new Size(_image.Width, _image.Height);
            int counter = 0;
            do
            {
                if (counter == zoomLevel)
                {
                    break;
                }
                size.Width = (int)(size.Width * ScaleDownPercentage.Factor);
                size.Height = (int)(size.Height * ScaleDownPercentage.Factor);
                counter += 1;
            } while (true);
            return size;
        }

        private List<int> RequiredZoomLevelIndexes()
        {
            var levels = new List<int>();
            int zoomLevelCount = 0;
            var imageSize = new Size(_image.Width, _image.Height);
            do
            {
                imageSize.Width = (int)(imageSize.Width * ScaleDownPercentage.Factor);
                imageSize.Height = (int)(imageSize.Height * ScaleDownPercentage.Factor);
                if (imageSize.Width < MinimumScaleSize.Width && imageSize.Height < MinimumScaleSize.Height)
                {
                    break;
                }
                levels.Add(zoomLevelCount);
                zoomLevelCount += 1;
            } while (true);
            return levels;
        }

        private string GetMetadataFilePath() => Path.Combine(_options.OutputDirectory, "meta.json");

        private string GetTileOutputPath(int zoomLevel, int row, int col)
        {
            string extension = Path.GetExtension(_options.File);
            if (string.IsNullOrWhiteSpace(extension))
            {
                extension = _image.Extension;
            }
            if (extension.StartsWith('.'))
            {
                extension = extension[1..];
            }
            return Path.Combine(_options.OutputDirectory, zoomLevel.ToString(), $"{row}x{col}.{extension}");
        }

        private string GetZoomLevelDirectoryPath(int zoomLevel) => Path.Combine(_options.OutputDirectory, zoomLevel.ToString());

        private void CreateZoomLevelDirectory(int zoomLevel)
        {
            CreateDirectory(GetZoomLevelDirectoryPath(zoomLevel));
        }

        private void CreateDirectory(string directory)
        {
            if (Directory.Exists(directory))
            {
                if (ActionOnExistingDirectory == ActionOnExistingDirectory.Delete)
                {
                    DeleteDirectoryWithContent(directory);
                }
                else if (ActionOnExistingDirectory == ActionOnExistingDirectory.Skip)
                {
                    return;
                }
            }
            Directory.CreateDirectory(directory);
        }
    }
}