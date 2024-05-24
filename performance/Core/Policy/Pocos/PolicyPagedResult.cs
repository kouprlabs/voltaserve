namespace Defyle.Core.Policy.Pocos
{
  using System.Collections.Generic;
  using Models;

  public class PolicyPagedResult
	{
    public IEnumerable<Policy> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}