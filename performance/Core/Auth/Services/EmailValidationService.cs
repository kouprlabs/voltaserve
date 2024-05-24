namespace Defyle.Core.Auth.Services
{
	public class EmailValidationService
	{
		public static bool IsValid(string email)
		{
			if (string.IsNullOrWhiteSpace(email))
			{
				return false;
			}

			try
			{
				// ReSharper disable once ObjectCreationAsStatement
				new System.Net.Mail.MailAddress(email);

				return true;
			}
			catch
			{
				return false;
			}
		}
	}
}