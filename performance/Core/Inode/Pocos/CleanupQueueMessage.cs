namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;

  public class CleanupQueueMessage
	{
    public string Id { get; set; }

    public IEnumerable<string> Paths { get; set; }
	}
}