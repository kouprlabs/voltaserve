namespace Defyle.Core.Inode.Pocos
{
  using System.Collections.Generic;
  using System.Linq;
  using Models;

  public class InodePagedResult
	{
		public InodePagedResult(IEnumerable<InodeFacet> data, long totalPages, long totalElements, long page, long size)
		{
			Data = data.Select(e => e);
			TotalPages = totalPages;
			TotalElements = totalElements;
			Page = page;
			Size = size;
		}

		public IEnumerable<InodeFacet> Data { get; }

		public long TotalPages { get; }

		public long TotalElements { get; }

		public long Page { get; }

		public long Size { get; }
	}
}