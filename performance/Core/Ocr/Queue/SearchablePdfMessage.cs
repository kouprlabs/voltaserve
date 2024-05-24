namespace Defyle.Core.Ocr.Queue
{
  using Newtonsoft.Json.Linq;

  public class SearchablePdfMessage
	{
    public string Id { get; set; }
    
		public string WorkspaceId { get; set; }

		public string InodeId { get; set; }

		public JObject Payload { get; set; }

		public void SetPayload<T>(T payload)
		{
			Payload = JObject.FromObject(payload);
		}

		public T GetPayload<T>()
		{
			return Payload.ToObject<T>();
		}
    
    public class SearchablePdfPayload
    {
      public string File { get; set; }

      public string OutputDirectory { get; set; }
    }
	}
}