namespace Voltaserve.Mosaic.Services
{
    using System;
    using System.IO;
    using Voltaserve.Mosaic.Infra;
    using Microsoft.AspNetCore.StaticFiles;
    using Models;
    using Minio;
    using Minio.DataModel.Args;
    using System.Threading.Tasks;
    using System.Text.Json;
    using Minio.Exceptions;
    using System.Collections.Generic;

    public class ResourceNotFoundException(string message) : Exception(message) { }


    public class MosaicService(IMinioClient _minioClient)
    {
        private readonly FileExtensionContentTypeProvider _fileExtensionContentTypeProvider = new();

        public async Task<Metadata> CreateAsync(string path, string s3Key, string s3Bucket)
        {
            var id = Ids.New();
            var outputDirectory = Path.Combine(Path.GetTempPath(), id);
            try
            {
                var metadata = new MosaicBuilder(new TilesBuilterOptions
                {
                    File = path,
                    OutputDirectory = outputDirectory,
                }).Build();

                var files = Directory.GetFiles(outputDirectory, "*.*", SearchOption.AllDirectories);
                foreach (var file in files)
                {
                    using var stream = new FileStream(file, FileMode.Open, FileAccess.Read);
                    await _minioClient.PutObjectAsync(new PutObjectArgs()
                        .WithBucket(s3Bucket)
                        .WithObject(Path.Combine(s3Key, "mosaic", Path.GetRelativePath(outputDirectory, file)))
                        .WithStreamData(stream)
                        .WithObjectSize(new FileInfo(file).Length)
                        .WithContentType(_fileExtensionContentTypeProvider.Mappings[Path.GetExtension(file)]));
                }
                return metadata;
            }
            catch
            {
                throw;
            }
            finally
            {
                if (Directory.Exists(outputDirectory))
                {
                    Directory.Delete(outputDirectory, true);
                }
            }
        }

        public async Task DeleteAsync(string s3Bucket, string s3Key)
        {
            try
            {
                var tcs = new TaskCompletionSource<bool>();
                var keysToDelete = new List<string>();
                _minioClient.ListObjectsAsync(new ListObjectsArgs()
                        .WithBucket(s3Bucket)
                        .WithPrefix(Path.Combine(s3Key, "mosaic"))
                        .WithRecursive(true))
                    .Subscribe(
                        item => keysToDelete.Add(item.Key),
                        () => tcs.SetResult(true));
                await tcs.Task;

                foreach (var key in keysToDelete)
                {
                    await _minioClient.RemoveObjectAsync(new RemoveObjectArgs()
                        .WithBucket(s3Bucket)
                        .WithObject(key));
                }
            }
            catch (MinioException)
            {
                throw new ResourceNotFoundException(null);
            }
        }

        public async Task<(Stream stream, string contentType)> GetTileStreamAsync(string s3Bucket, string s3Key, int zoomLevel, int row, int col, string ext)
        {
            try
            {
                var memoryStream = new MemoryStream();
                await _minioClient.GetObjectAsync(new GetObjectArgs()
                    .WithBucket(s3Bucket)
                    .WithObject(Path.Combine(s3Key, "mosaic", zoomLevel.ToString(), $"{row}x{col}.{ext}"))
                    .WithCallbackStream((stream) => stream.CopyTo(memoryStream)));
                memoryStream.Seek(0, SeekOrigin.Begin);
                return (memoryStream, _fileExtensionContentTypeProvider.Mappings[$".{ext}"]);
            }
            catch (MinioException)
            {
                throw new ResourceNotFoundException(null);
            }
        }

        public async Task<Metadata> GetMetadataAsync(string s3Bucket, string s3Key)
        {
            try
            {
                var memoryStream = new MemoryStream();
                await _minioClient.GetObjectAsync(new GetObjectArgs()
                    .WithBucket(s3Bucket)
                    .WithObject(Path.Combine(s3Key, "mosaic", "meta.json"))
                    .WithCallbackStream((stream) => stream.CopyTo(memoryStream)));
                memoryStream.Seek(0, SeekOrigin.Begin);
                return await JsonSerializer.DeserializeAsync<Metadata>(memoryStream);
            }
            catch (MinioException)
            {
                throw new ResourceNotFoundException(null);
            }
        }
    }
}