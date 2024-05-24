namespace Defyle.Core.Preview.Queue
{
  using Newtonsoft.Json.Linq;

  public class PreviewMessage
	{
    public string Id { get; set; }
    
		public string WorkspaceId { get; set; }

		public string InodeId { get; set; }

		public JObject Payload { get; set; }

    public bool S3Enabled { get; set; }

		public void SetPayload<T>(T payload)
		{
			Payload = JObject.FromObject(payload);
		}

		public T GetPayload<T>()
		{
			return Payload.ToObject<T>();
		}
    
    public class DocumentPreviewPayload
    {
      public string File { get; set; }

      public string OutputDirectory { get; set; }
    }
    
    public class ImagePreviewPayload
    {
      public string File { get; set; }

      public string OutputDirectory { get; set; }

      public int PreviewSize { get; set; }

      public string Extension { get; set; }
    }
    
    public class TileMapPayload
    {
      public string File { get; set; }

      public string OutputDirectory { get; set; }

      public string Extension { get; set; }
    }
	}
}