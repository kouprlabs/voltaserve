namespace Defyle.WebApi.Infrastructure.Dtos
{
	public class VersionDto
	{
		public int Major { get; set; }

		public int Minor { get; set; }

		public int Patch { get; set; }

		public string PreRelease { get; set; }

		public string Build { get; set; }
	}
}