namespace Defyle.WebApi.Inode.Requests
{
  public class InodeSearchRequest
	{
		public string SearchCriteria { get; set; }

    public string ParentId { get; set; }
    
    public bool? IncludeFiles { get; set; }

    public bool? IncludeDirectories { get; set; }

    public long? CreatedAtFrom { get; set; }

		public long? CreatedAtTo { get; set; }

		public long? UpdatedAtFrom { get; set; }

		public long? UpdatedAtTo { get; set; }
	}
}