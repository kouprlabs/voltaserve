namespace Defyle.Core.Auth.Extensions
{
  using System;
  using Models;

  public static class AuthenticationTypeExtensions
	{
		private const string Ldap = "ldap";
		private const string Local = "local";

		public static string ToAuthenticationTypeString(this AuthenticationType type)
		{
			switch (type)
			{
				case AuthenticationType.Ldap:

					return Ldap;

				case AuthenticationType.Local:

					return Local;

				default:

					throw new Exception($"Unknown authentication type '{type.ToString()}'");
			}
		}

		public static AuthenticationType ToAuthenticationType(this string type)
		{
			string normalized = type.NormalizedAuthenticationTypeString();

			switch (normalized)
			{
				case Ldap:

					return AuthenticationType.Ldap;

				case Local:

					return AuthenticationType.Local;

				default:

					throw new Exception($"Invalid authentication type: '{type}'");
			}
		}

		public static string NormalizedAuthenticationTypeString(this string type)
		{
			if (string.IsNullOrWhiteSpace(type))
			{
				throw new Exception($"Invalid authentication type: '{type}'");
			}

			string normalized = type.Trim().ToLowerInvariant();

			switch (normalized)
			{
				case Ldap:
				case Local:

					return normalized;

				default:

					throw new Exception($"Invalid authentication type: '{type}'");
			}
		}

		public static bool IsLdap(this AuthenticationType authenticationType)
		{
			if (authenticationType == AuthenticationType.Ldap)
			{
				return true;
			}

			return false;
		}
	}
}