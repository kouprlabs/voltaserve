namespace Defyle.Core.Auth.Services
{
  using System;
  using System.Collections.Generic;
  using System.Drawing;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Exceptions;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Models;
  using Preview.Services;
  using Proto;

  public class AccountService : QueueService
	{
		private readonly Size _profileImageSize = new Size(230, 230);
		private readonly CoreSettings _coreSettings;
    private readonly PasswordService _passwordService;
    private readonly EmailQueueService _emailQueueService;
    private readonly UserServiceProto.UserServiceProtoClient _protoClient;
    private readonly IMapper _mapper;
    private readonly UserService _userService;

    public AccountService(
			CoreSettings coreSettings,
      PasswordService passwordService,
      EmailQueueService emailQueueService,
      UserServiceProto.UserServiceProtoClient protoClient,
      IMapper mapper,
      UserService userService)
			: base(coreSettings.MessageBroker)
		{
			_coreSettings = coreSettings;
      _passwordService = passwordService;
      _emailQueueService = emailQueueService;
      _protoClient = protoClient;
      _mapper = mapper;
      _userService = userService;
    }

    public async Task<User> CreateLocalAsync(User user)
		{
      user.FullName = user.FullName.Trim();
      user.Email = user.Email.ToLowerInvariant().Trim();

			List<Error> passwordValidationErrors = _passwordService.ValidatePassword(user.PasswordHash);
      if (passwordValidationErrors.Any())
			{
				throw new PasswordValidationException().WithErrors(passwordValidationErrors);
			}

			List<Error> emailValidationErrors = await _userService.ValidateEmailAsync(user.Email);
      if (emailValidationErrors.Any())
			{
				throw new EmailValidationException().WithErrors(emailValidationErrors);
			}

			if (!Base64ImageService.IsValidBase64Image(user.Image))
			{
				if (_coreSettings.GravatarIntegration)
				{
          user.Image = Base64ImageService.GenerateGravatar(user.Email, _profileImageSize.Width);
				}
			}

      string emailConfirmationToken = Guid.NewGuid().ToString().Replace("-", string.Empty);
      
      User systemUser = await _userService.FindSystemUserAsync();

      var model = new User();
      model.Email = user.Email.Trim().ToLowerInvariant();
      model.Username = user.Email.Trim().ToLowerInvariant();
      model.FullName = user.FullName.Trim();
      model.Image = user.Image;
      model.PasswordHash = _passwordService.HashPassword(user.PasswordHash);
      model.IsEmailConfirmed = false;
      model.EmailConfirmationToken = emailConfirmationToken;
      model.CreateTime = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();

      var proto = await _protoClient.InsertAsync(new UserInsertRequestProto
      {
        UserId = systemUser.Id,
        User = _mapper.Map<UserProto>(model)
      });
      var response = _mapper.Map<User>(proto);

			EmailMessage message = new EmailMessage();
      message.ToEmail = response.Email;
      message.ToName = model.FullName;
      message.Subject = _coreSettings.ConfirmationEmail.Subject;
      message.HtmlContent = ReplaceVariables(_coreSettings.ConfirmationEmail.GetHtmlContent(),
        emailConfirmationToken, _coreSettings.WebClientUrl);
      message.PlainTextContent = ReplaceVariables(_coreSettings.ConfirmationEmail.GetPlainTextContent(),
        emailConfirmationToken, _coreSettings.WebClientUrl);;

      _emailQueueService.SendEmailMessage(message);

      return response;
		}

    public async Task<User> ConfirmEmailAsync(string token)
		{
      User systemUser = await _userService.FindSystemUserAsync();
      User user = await _userService.FindByEmailConfirmationTokenAsync(token, systemUser);
      
      user.IsEmailConfirmed = true;
      user.EmailConfirmationToken = null;

      await _protoClient.UpdateAsync(new UserUpdateRequestProto
      {
        UserId = systemUser.Id,
        User = _mapper.Map<UserProto>(user)
      });

			return user;
		}

    public async Task SendResetPasswordEmailAsync(string email)
		{
      email = email.Trim();
      
      string resetPasswordToken = Guid.NewGuid().ToString().Replace("-", string.Empty);
      
      User systemUser = await _userService.FindSystemUserAsync();
      User user = await _userService.FindByEmailAsync(email, systemUser);
      user.ResetPasswordToken = resetPasswordToken;

      await _protoClient.UpdateAsync(new UserUpdateRequestProto
      {
        UserId = systemUser.Id,
        User = _mapper.Map<UserProto>(user)
      });
      
			EmailMessage message = new EmailMessage();
      message.ToEmail = user.Email;
      message.ToName = user.FullName;
      message.Subject = _coreSettings.ResetPasswordEmail.Subject;
      message.HtmlContent = ReplaceVariables(_coreSettings.ResetPasswordEmail.GetHtmlContent(),
        resetPasswordToken, _coreSettings.WebClientUrl);
      message.PlainTextContent = ReplaceVariables(_coreSettings.ResetPasswordEmail.GetPlainTextContent(),
        resetPasswordToken, _coreSettings.WebClientUrl);;

      _emailQueueService.SendEmailMessage(message);
		}

		public async Task ResetPasswordAsync(string token, string newPassword)
		{
      User systemUser = await _userService.FindSystemUserAsync();
      User user = await _userService.FindByResetPasswordTokenAsync(token, systemUser);
      
			List<Error> passwordValidationErrors = _passwordService.ValidatePassword(newPassword).ToList();

			if (passwordValidationErrors.Any())
			{
				throw new PasswordValidationException().WithErrors(passwordValidationErrors);
			}
			else
      {
        user.PasswordHash = _passwordService.HashPassword(newPassword);
        user.ResetPasswordToken = null;
        
        await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = systemUser.Id,
          User = _mapper.Map<UserProto>(user)
        });
      }
		}

    private static string ReplaceVariables(string text, string token, string webClientUrl)
    {
      return text
        .Replace("#TOKEN", token)
        .Replace("#WEB_CLIENT_URL", webClientUrl);
    }
  }
}