namespace Defyle.Core.Infrastructure.Poco
{
  using Newtonsoft.Json;

  public class S3Settings
  {
    [JsonProperty("awsAccessKeyId")]
    public string AwsAccessKeyId { get; set; }

    [JsonProperty("awsSecretKey")]
    public string AwsSecretKey { get; set; }
    
    [JsonProperty("awsRegion")]
    public string AwsRegion { get; set; }

    [JsonProperty("s3Bucket")]
    public string S3Bucket { get; set; }
  }
}