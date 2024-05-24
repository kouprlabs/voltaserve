namespace Defyle.WebApi.Policy.Dtos
{
  using System.Collections.Generic;

  public class PolicyPagedResultDto
	{
		public IEnumerable<PolicyDto> Data { get; set; }

		public long TotalPages { get; set; }

		public long TotalElements { get; set; }

		public long Page { get; set; }

		public long Size { get; set; }
	}
}