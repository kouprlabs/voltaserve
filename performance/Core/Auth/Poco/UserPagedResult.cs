namespace Defyle.Core.Auth.Poco
{
  using System.Collections.Generic;
  using Models;

  public class UserPagedResult
	{
    public IEnumerable<User> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}