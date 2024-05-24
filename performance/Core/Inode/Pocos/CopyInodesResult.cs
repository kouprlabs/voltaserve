namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using Infrastructure.Poco;

  public class CopyInodesResult
	{
		public List<string> Succeeded { get; set; } = new List<string>();
    
		public List<string> Failed { get; set; } = new List<string>();
    
		public List<string> Created { get; set; } = new List<string>();
    
    public List<Error> Errors { get; set; } = new List<Error>();
	}
}