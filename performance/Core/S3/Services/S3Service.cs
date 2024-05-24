namespace Defyle.Core.S3.Services
{
  using System.IO;
  using System.Threading.Tasks;
  using Amazon;
  using Amazon.Runtime;
  using Amazon.S3;
  using Amazon.S3.Transfer;
  using Infrastructure.Poco;

  public class S3Service
  {
    private readonly S3Settings _settings;

    public S3Service(S3Settings settings)
    {
      _settings = settings;
    }
    
    public async Task UploadAsync(string localPath, string s3Key)
    {
      using var client = CreateClient();
      await using var memoryStream = new MemoryStream();
      await using (FileStream fileStream = File.Open(localPath, FileMode.Open))
      {
        fileStream.CopyTo(memoryStream);
      }

      var request = new TransferUtilityUploadRequest
      {
        InputStream = memoryStream,
        Key = s3Key,
        BucketName = _settings.S3Bucket
      };

      var transferUtility = new TransferUtility(client);
      await transferUtility.UploadAsync(request);
    }

    public async Task DownloadAsync(string localPath, string s3Key)
    {
      using var client = CreateClient();
      var transferUtility = new TransferUtility(client);
      await transferUtility.DownloadAsync(localPath, _settings.S3Bucket, s3Key);
    }

    private AmazonS3Client CreateClient() =>
      new AmazonS3Client(new BasicAWSCredentials(_settings.AwsAccessKeyId, _settings.AwsSecretKey),
        RegionEndpoint.GetBySystemName(_settings.AwsRegion));
  }
}