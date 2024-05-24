namespace Defyle.WebApi.Inode.Dtos
{
  using System.Collections.Generic;

  public class InodePagedResultDto
	{
		public IEnumerable<InodeFacetDto> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}