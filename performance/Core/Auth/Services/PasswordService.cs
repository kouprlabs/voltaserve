namespace Defyle.Core.Auth.Services
{
  using System;
  using System.Collections.Generic;
  using Infrastructure.Poco;
  using Microsoft.AspNetCore.Identity;
  using Models;

  public class PasswordService
	{
		private readonly CoreSettings _coreSettings;
		private readonly IPasswordHasher<User> _passwordHasher;
		private readonly PasswordValidator<User> _passwordValidator;

		public PasswordService(CoreSettings coreSettings)
		{
			_coreSettings = coreSettings;
			_passwordHasher = new PasswordHasher<User>();
			_passwordValidator = new PasswordValidator<User>();
		}

		public bool VerifyHashedPassword(string passwordHash, string password)
		{
			PasswordVerificationResult result = _passwordHasher.VerifyHashedPassword(null, passwordHash, password);

			if (result == PasswordVerificationResult.Failed)
			{
				return false;
			}

			return true;
		}

		public string HashPassword(string password)
		{
			return _passwordHasher.HashPassword(null, password);
		}

		public List<Error> ValidatePassword(string password)
		{
			List<Error> errors = new List<Error>();

			if (string.IsNullOrEmpty(password))
			{
				errors.Add(Error.NullOrEmptyPasswordError);
        return errors;
			}

			if (password.Length < _coreSettings.Password.RequiredLength)
			{
				errors.Add(new Error(
					"InvalidPasswordLength",
					"Invalid password length."));
			}

			if (_coreSettings.Password.Digits)
			{
				if (!Fulfills(password, c => _passwordValidator.IsDigit(c)))
				{
					errors.Add(new Error(
						"PasswordRequiresDigit",
						"Password requires at least one digit character."));
				}
			}

			if (_coreSettings.Password.Lowercase)
			{
				if (!Fulfills(password, c => _passwordValidator.IsLower(c)))
				{
					errors.Add(new Error(
						"PasswordRequiresLowercase",
						"Password requires at least one lowercase character."));
				}
			}

			if (_coreSettings.Password.Uppercase)
			{
				if (!Fulfills(password, c => _passwordValidator.IsUpper(c)))
				{
					errors.Add(new Error(
						"PasswordRequiresUppercase",
						"Password requires at least one uppercase character."));
				}
			}

			if (_coreSettings.Password.NonAlphanumeric)
			{
				if (!Fulfills(password, c => !_passwordValidator.IsLetterOrDigit(c)))
				{
					errors.Add(new Error(
						"PasswordRequiresNonAlphanumeric",
						"Password requires at least one non-alphanumeric character."));
				}
			}

			if (_coreSettings.Password.Unique)
			{
				if (IsUnique(password))
				{
					errors.Add(new Error(
						"PasswordRequiresUniqueCharacters",
						"Password requires unique characters."));
				}
			}

			return errors;
		}

		private static bool Fulfills(string value, Func<char, bool> predicate)
		{
			foreach (char c in value)
			{
				if (predicate(c))
				{
					return true;
				}
			}

			return false;
		}

		private static bool IsUnique(string value)
		{
			var set = new HashSet<char>();

			foreach (var c in value)
			{
				if (!set.Add(c))
					return false;
			}

			return true;
		}
	}
}