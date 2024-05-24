namespace Defyle.Core.Preview.Services
{
  using System;
  using System.Net;
  using System.Security.Cryptography;
  using System.Text;

  public class Base64ImageService
	{
		public const string Base64ImagePrefix = "data:image/png;base64,";

		public static bool IsValidBase64Image(string value)
		{
			return !string.IsNullOrWhiteSpace(value) && value.StartsWith(Base64ImagePrefix);
		}

		public static string GenerateGravatar(string email, int size)
		{
			try
			{
				byte[] data = MD5.Create().ComputeHash(Encoding.Default.GetBytes(email));

				StringBuilder emailHashBuilder = new StringBuilder();

				foreach (var value in data)
				{
					emailHashBuilder.Append(value.ToString("x2"));
				}

				string image = DownloadImageToBase64($"http://gravatar.com/avatar/{emailHashBuilder}?s={size}&d=404");

				if (!IsValidBase64Image(image))
				{
					return null;
				}

				return image;
			}
			catch
			{
				return null;
			}
		}

		private static string DownloadImageToBase64(string url)
    {
      using WebClient client = new WebClient();
      byte[] bytes = client.DownloadData(new Uri(url));

      return "data:image/png;base64," + Convert.ToBase64String(bytes);
    }
	}
}