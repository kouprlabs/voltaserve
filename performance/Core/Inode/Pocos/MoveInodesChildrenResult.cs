namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class MoveInodesChildrenResult
	{
		public List<string> PreviousParentIds { get; set; } = new List<string>();
    
		public string CurrentParentId { get; set; }
    
		public List<string> Succeeded { get; set; } = new List<string>();
    
		public List<string> Failed { get; set; } = new List<string>();
    
    public List<Error> Errors { get; set; } = new List<Error>();
	}
}