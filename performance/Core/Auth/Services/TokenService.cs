namespace Defyle.Core.Auth.Services
{
  using System;
  using System.Collections.Generic;
  using System.IdentityModel.Tokens.Jwt;
  using System.Security.Claims;
  using System.Text;
  using System.Threading.Tasks;
  using AutoMapper;
  using Exceptions;
  using Extensions;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Microsoft.IdentityModel.Tokens;
  using Models;
  using Newtonsoft.Json;
  using Novell.Directory.Ldap;
  using Poco;
  using Policy.Services;
  using Proto;

  public class TokenService
	{
		private readonly CoreSettings _coreSettings;
    private readonly PasswordService _passwordService;
    private readonly LdapService _ldapService;
    private readonly UserService _userService;
    private readonly UserServiceProto.UserServiceProtoClient _protoClient;
    private readonly PolicyService _policyService;
    private readonly IMapper _mapper;
    private readonly SymmetricSecurityKey _issuerSigninKey;
		private readonly JwtSecurityTokenHandler _jwtSecurityTokenHandler = new JwtSecurityTokenHandler();
		private readonly LdapConnection _ldapConnection;

		public TokenService(
			CoreSettings coreSettings,
      PasswordService passwordService,
      LdapService ldapService,
			UserService userService,
      UserServiceProto.UserServiceProtoClient protoClient,
      PolicyService policyService,
      IMapper mapper)
		{
			_coreSettings = coreSettings;
      _passwordService = passwordService;
      _ldapService = ldapService;
      _userService = userService;
      _protoClient = protoClient;
      _policyService = policyService;
      _mapper = mapper;
      _issuerSigninKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_coreSettings.SecurityKey));

			TokenValidationParameters = new TokenValidationParameters
			{
				ValidIssuer = _coreSettings.Token.TokenIssuer,
				ValidAudience = _coreSettings.Token.TokenAudience,
				ValidateIssuer = true,
				ValidateAudience = true,
				ValidateLifetime = true,
				ValidateIssuerSigningKey = true,
				IssuerSigningKey = _issuerSigninKey
			};

			_ldapConnection = new LdapConnection();
		}

		public TokenValidationParameters TokenValidationParameters { get; }

		public async Task<Token> ExchangeAsync(TokenExchangeOptions options)
		{
			if (options.GrantType == "password")
			{
				User user;

				if (_coreSettings.AuthenticationType.ToAuthenticationType() == AuthenticationType.Local)
				{
					user = await LocalLoginAsync(options.Username, options.Password);
				}
				else if (_coreSettings.AuthenticationType.ToAuthenticationType() == AuthenticationType.Ldap)
				{
					user = await LdapLoginAsync(options.Username, options.Password);
				}
				else
				{
					throw new Exception("Invalid authentication type");
				}

				RefreshToken refreshToken = await CreateAndPersistRefreshTokenAsync(user);
				JwtSecurityToken accessToken = CreateAccessToken(user);

				Token token = new Token();
				token.AccessToken = _jwtSecurityTokenHandler.WriteToken(accessToken);
				token.TokenType = "Bearer";
				token.ExpiresIn = _coreSettings.Token.AccessTokenLifetime;
				token.RefreshToken = refreshToken.Value;

        User systemUser = await _userService.FindSystemUserAsync();
        await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = systemUser.Id,
          User = _mapper.Map<UserProto>(user)
        });

				return token;
			}
			else if (options.GrantType == "refresh_token")
			{
				User user;
				try
				{
          User systemUser = await _userService.FindSystemUserAsync();
					user = await _userService.FindByRefreshTokenAsync(options.RefreshToken, systemUser);
				}
				catch (ResourceNotFoundException)
				{
					throw new AuthenticationException();
				}

				if (DateTimeOffset.UtcNow > DateTimeOffset.FromUnixTimeMilliseconds(user.RefreshTokenValidTo))
				{
					throw new AuthenticationException();
				}

				RefreshToken newRefreshToken = await CreateAndPersistRefreshTokenAsync(user);
				JwtSecurityToken accessToken = CreateAccessToken(user);

				Token token = new Token();
				token.AccessToken = _jwtSecurityTokenHandler.WriteToken(accessToken);
				token.TokenType = "Bearer";
				token.ExpiresIn = _coreSettings.Token.AccessTokenLifetime;
				token.RefreshToken = newRefreshToken.Value;

				return token;
			}
			else
			{
				throw new ArgumentException("Invalid grant");
			}
		}

		private JwtSecurityToken CreateAccessToken(User user)
		{
			var claims = new List<Claim>();

			claims.Add(new Claim(JwtRegisteredClaimNames.Sub, user.Id));

			claims.Add(new Claim(JwtRegisteredClaimNames.UniqueName, user.Username));

			if (!string.IsNullOrWhiteSpace(user.Email))
			{
				claims.Add(new Claim(JwtRegisteredClaimNames.Email, user.Email));
			}

      User systemUser = _userService.FindSystemUserAsync().Result;

      IEnumerable<string> systemModelPermissions = _policyService.GetModelPermissionsForUser(user.Id, "system", systemUser);
      
			claims.Add(new Claim("permissions", JsonConvert.SerializeObject(systemModelPermissions)));

      return new JwtSecurityToken(
				_coreSettings.Token.TokenIssuer,
				_coreSettings.Token.TokenAudience,
				expires: DateTime.UtcNow.AddSeconds(_coreSettings.Token.AccessTokenLifetime),
				claims: claims,
				signingCredentials: new SigningCredentials(_issuerSigninKey, SecurityAlgorithms.HmacSha256));
		}

		private async Task<RefreshToken> CreateAndPersistRefreshTokenAsync(User user)
		{
      RefreshToken refreshToken = new RefreshToken();
      refreshToken.Value = Guid.NewGuid().ToString().Replace("-", string.Empty);
      refreshToken.ValidTo = DateTimeOffset.UtcNow.AddSeconds(_coreSettings.Token.RefreshTokenLifetime).ToUnixTimeMilliseconds();

      user.RefreshTokenValue = refreshToken.Value;
      user.RefreshTokenValidTo = refreshToken.ValidTo;

      User systemUser = await _userService.FindSystemUserAsync();
      await _protoClient.UpdateAsync(new UserUpdateRequestProto
      {
        UserId = systemUser.Id,
        User = _mapper.Map<UserProto>(user)
      });
      
			return refreshToken;
		}

		private async Task<User> LocalLoginAsync(string username, string password)
		{
      User systemUser = await _userService.FindSystemUserAsync();
			User user = await _userService.FindByUsernameAsync(username, systemUser);
      if (user == null)
			{
				throw new AuthenticationException().WithError(Error.InvalidEmailOrPasswordError);
			}

			if (!user.IsEmailConfirmed)
			{
				throw new EmailNotConfirmedException().WithError(Error.EmailNotConfirmedError);
			}

			if (!_passwordService.VerifyHashedPassword(user.PasswordHash, password))
			{
				throw new AuthenticationException().WithError(Error.InvalidEmailOrPasswordError);
			}

			return user;
		}

		private async Task<User> LdapLoginAsync(string username, string password)
		{
			try
			{
				_ldapConnection.Connect(_coreSettings.Ldap.Host, _coreSettings.Ldap.Port);
			}
			catch
			{
				throw new AuthenticationException().WithError(Error.LdapConnectionError);
			}

			try
			{
				_ldapConnection.Bind(_coreSettings.Ldap.BindDn, _coreSettings.Ldap.BindPassword);
			}
			catch
			{
				throw new AuthenticationException().WithError(Error.LdapBindError);
			}

			try
			{
				string searchFilter = string.Format(_coreSettings.Ldap.SearchFilter, username);

				var result = _ldapConnection.Search(
					_coreSettings.Ldap.SearchBase,
					LdapConnection.ScopeSub,
					searchFilter,
					new[] {"memberOf", "displayName", "sAMAccountName", "mail"},
					false
				);

				LdapEntry ldapEntry = result.Next();

				if (ldapEntry != null)
				{
					_ldapConnection.Bind(ldapEntry.Dn, password);

					if (_ldapConnection.Bound)
					{
						return await _ldapService.CreateFromLdapEntryAsync(ldapEntry);
					}
				}
			}
			catch
			{
				throw new AuthenticationException().WithError(Error.LdapAuthenticationError);
			}

			_ldapConnection.Disconnect();

			return null;
		}
	}
}