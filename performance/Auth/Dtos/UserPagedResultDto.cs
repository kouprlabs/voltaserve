namespace Defyle.WebApi.Auth.Dtos
{
  using System.Collections.Generic;

  public class UserPagedResultDto
	{
		public IEnumerable<UserDto> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}