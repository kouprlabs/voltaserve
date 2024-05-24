namespace Defyle.Core.Storage.Models
{
	public class StorageUsage
	{
		public StorageUsage()
		{
		}

		public StorageUsage(long bytes, long maxBytes)
		{
			Bytes = bytes;
			MaxBytes = maxBytes;
      if (maxBytes != 0)
      {
        Percentage = (int) (Bytes * 100 / maxBytes);
      }
    }

		public long Bytes { get; set; }

		public long MaxBytes { get; set; }

		public int Percentage { get; set; }
	}
}