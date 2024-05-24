namespace Defyle.Core.Streaming.Services
{
  using System;
  using System.IO;
  using System.Threading.Tasks;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using S3.Services;
  using Storage.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public abstract class StreamService : IStreamService
  {
    private readonly CoreSettings _coreSettings;

    public StreamService(CoreSettings coreSettings)
    {
      _coreSettings = coreSettings;
    }
    
    public async Task<Stream> GetStreamAsync(Workspace workspace, File file, string cipherPassword)
    {
      string localPath;
      if (_coreSettings.S3Enabled)
      {
        localPath = Path.Combine(_coreSettings.TempDirectory, Guid.NewGuid().ToString());
        
        S3Service s3Service = new S3Service(_coreSettings.S3);
        await s3Service.DownloadAsync(localPath, await GetS3KeyAsync(workspace, file));
      }
      else
      {
        localPath = await GetLocalPathAsync(workspace, file);
      }
      
      if (!System.IO.File.Exists(localPath))
      {
        throw new ResourceNotFoundException().WithError(Error.PhysicalFileNotFoundError);
      }

      if (!workspace.Encrypted)
      {
        return new FileStream(localPath, FileMode.Open);
      }
      
      if (string.IsNullOrWhiteSpace(cipherPassword))
      {
        throw new InternalServerErrorException().WithError(Error.WorkspacePasswordRequiredError);
      }

      byte[] key = EncryptionService.EncryptionKey(workspace, cipherPassword);

      MemoryStream outputStream = new MemoryStream();
      await using (FileStream inputStream = new FileStream(localPath, FileMode.Open))
      {
        EncryptionService.DecryptStreamWithSalt(inputStream, outputStream, key);
      }
      outputStream.Position = 0;

      return outputStream;
    }

    public abstract Task<string> GetLocalPathAsync(Workspace workspace, File file);
    
    public abstract Task<string> GetLocalDirectoryAsync(Workspace workspace, File file);

    public abstract Task<string> GetS3KeyAsync(Workspace workspace, File file);
  }
}