namespace Defyle.WebApi.Role.Dtos
{
  using System.Collections.Generic;

  public class RolePagedResultDto
	{
		public IEnumerable<RoleDto> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}