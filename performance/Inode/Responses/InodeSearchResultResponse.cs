namespace Defyle.WebApi.Inode.Responses
{
  using Dtos;
  using Requests;

  public class InodeSearchResultResponse
	{
		public InodePagedResultDto Result { get; set; }

		public InodeSearchRequest Request { get; set; }
	}
}