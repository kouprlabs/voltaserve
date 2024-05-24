namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class DeleteInodesResult
	{
		public List<string> AffectedParentIds { get; set; } = new List<string>();
    
		public List<string> Succeeded { get; set; } = new List<string>();
    
		public List<string> Failed { get; set; } = new List<string>();
    
    public List<Error> Errors { get; set; } = new List<Error>();
	}
}