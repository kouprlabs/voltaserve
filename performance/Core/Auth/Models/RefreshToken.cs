namespace Defyle.Core.Auth.Models
{
	public class RefreshToken
	{
		public string Value { get; set; }

		public long ValidTo { get; set; }
	}
}