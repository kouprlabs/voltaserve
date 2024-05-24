namespace Defyle.Core.Role.Pocos
{
  using System.Collections.Generic;
  using Models;

  public class RolePagedResult
	{
    public IEnumerable<Role> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}