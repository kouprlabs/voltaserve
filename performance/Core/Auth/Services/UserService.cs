namespace Defyle.Core.Auth.Services
{
  using System;
  using System.Collections.Generic;
  using System.Drawing;
  using System.IO;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Exceptions;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Inode.Services;
  using Microsoft.AspNetCore.Identity;
  using Models;
  using Poco;
  using Preview.Models;
  using Preview.Services;
  using Proto;
  using Storage.Models;
  using Workspace.Models;
  using Workspace.Services;

  public class UserService
  {
    private readonly Size _profileImageSize = new Size(230, 230);
    private readonly CoreSettings _coreSettings;
    private readonly InodeEngine _inodeEngine;
    private readonly PasswordService _passwordService;
    private readonly UserServiceProto.UserServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;
    private readonly WorkspaceService _workspaceService;

    public UserService(
      CoreSettings coreSettings,
      InodeEngine inodeEngine,
      PasswordService passwordService,
      IMapper mapper,
      WorkspaceService workspaceService,
      UserServiceProto.UserServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator)
    {
      _coreSettings = coreSettings;
      _inodeEngine = inodeEngine;
      _passwordService = passwordService;
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
      _workspaceService = workspaceService;
    }

    public async Task<User> InsertInsecureAsync(User subject)
    {
      subject.FullName = subject.FullName.Trim();
      subject.Email = subject.Email.ToLowerInvariant().Trim();

      List<Error> passwordValidationErrors = _passwordService.ValidatePassword(subject.PasswordHash);
      if (passwordValidationErrors.Any())
      {
        throw new PasswordValidationException().WithErrors(passwordValidationErrors);
      }

      List<Error> emailValidationErrors = await ValidateEmailAsync(subject.Email);
      if (emailValidationErrors.Any())
      {
        throw new EmailValidationException().WithErrors(emailValidationErrors);
      }

      if (!Base64ImageService.IsValidBase64Image(subject.Image))
      {
        if (_coreSettings.GravatarIntegration)
        {
          subject.Image = Base64ImageService.GenerateGravatar(subject.Email, _profileImageSize.Width);
        }
      }

      var model = new User();
      model.Email = subject.Email.Trim().ToLowerInvariant();
      model.Username = subject.Username.Trim();
      model.FullName = subject.FullName.Trim();
      model.Image = subject.Image;
      model.PasswordHash = _passwordService.HashPassword(subject.PasswordHash);
      model.IsEmailConfirmed = true;
      model.IsSuperuser = subject.IsSuperuser;
      model.IsSystem = subject.IsSystem;
      model.CreateTime = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();

      try
      {
        var proto = await _protoClient.InsertInsecureAsync(_mapper.Map<UserProto>(model));
        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> InsertAsync(User subject, User user)
    {
      subject.FullName = subject.FullName.Trim();
      subject.Email = subject.Email.ToLowerInvariant().Trim();

      List<Error> passwordValidationErrors = _passwordService.ValidatePassword(subject.PasswordHash);
      if (passwordValidationErrors.Any())
      {
        throw new PasswordValidationException().WithErrors(passwordValidationErrors);
      }

      List<Error> emailValidationErrors = await ValidateEmailAsync(subject.Email);
      if (emailValidationErrors.Any())
      {
        throw new EmailValidationException().WithErrors(emailValidationErrors);
      }

      if (!Base64ImageService.IsValidBase64Image(subject.Image))
      {
        if (_coreSettings.GravatarIntegration)
        {
          subject.Image = Base64ImageService.GenerateGravatar(subject.Email, _profileImageSize.Width);
        }
      }

      var model = new User();
      model.Email = subject.Email.Trim().ToLowerInvariant();
      model.Username = subject.Username.Trim();
      model.FullName = subject.FullName.Trim();
      model.Image = subject.Image;
      model.PasswordHash = _passwordService.HashPassword(subject.PasswordHash);
      model.IsEmailConfirmed = true;
      model.IsSuperuser = subject.IsSuperuser;
      model.CreateTime = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();

      try
      {
        var proto = await _protoClient.InsertAsync(new UserInsertRequestProto
        {
          UserId = user.Id,
          User = _mapper.Map<UserProto>(model)
        });

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindAsync(string id)
    {
      User systemUser = await FindSystemUserAsync();
      var request = new UserFindRequestProto {UserId = systemUser.Id, Id = id};
      UserProto proto = await _protoClient.FindAsync(request);

      return _mapper.Map<User>(proto);
    }

    public async Task<User> FindAsync(string id, User user)
    {
      try
      {
        var request = new UserFindRequestProto {UserId = user.Id, Id = id};
        UserProto proto = await _protoClient.FindAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindByEmailAsync(string email, User user)
    {
      try
      {
        var request = new UserFindByFieldRequestProto {UserId = user.Id, Field = "email", Value = email};
        UserProto proto = await _protoClient.FindByFieldAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindByUsernameAsync(string username, User user)
    {
      try
      {
        var request = new UserFindByFieldRequestProto {UserId = user.Id, Field = "username", Value = username};
        UserProto proto = await _protoClient.FindByFieldAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindByRefreshTokenAsync(string refreshToken, User user)
    {
      try
      {
        var request = new UserFindByFieldRequestProto {UserId = user.Id, Field = "refresh_token_value", Value = refreshToken};
        UserProto proto = await _protoClient.FindByFieldAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindByEmailConfirmationTokenAsync(string emailConfirmationToken, User user)
    {
      try
      {
        var request = new UserFindByFieldRequestProto {UserId = user.Id, Field = "email_confirmation_token", Value = emailConfirmationToken};
        UserProto proto = await _protoClient.FindByFieldAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindByResetPasswordTokenAsync(string resetPasswordToken, User user)
    {
      try
      {
        var request = new UserFindByFieldRequestProto {UserId = user.Id, Field = "reset_password_token", Value = resetPasswordToken};
        UserProto proto = await _protoClient.FindByFieldAsync(request);

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<IEnumerable<User>> FindAllAsync(User user)
    {
      try
      {
        var request = new UserFindAllRequestProto {UserId = user.Id};
        UserFindAllResponseProto response = await _protoClient.FindAllAsync(request);

        return response.Users.Select(p => _mapper.Map<User>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<UserPagedResult> FindAllPagedAsync(User user, int page, int size)
    {
      try
      {
        var proto = await _protoClient.FindAllPagedAsync(new UserFindAllPagedRequestProto
        {
          UserId = user.Id,
          Page = page,
          Size = size
        });

        return _mapper.Map<UserPagedResult>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<StorageUsage> GetStorageUsageAsync(User user)
    {
      User systemUser = await FindSystemUserAsync();
      IEnumerable<Workspace> workspaces = await _workspaceService.FindAllAsync(user);
      long consumed = 0;
      long workspacesMaxStorageSum = 0;
      foreach (Workspace workspace in workspaces)
      {
        var rootNode = await _inodeEngine.FindRootAsync(workspace, systemUser);
        consumed += await _inodeEngine.GetStorageUsageAsync(rootNode, systemUser);
        workspacesMaxStorageSum += workspace.StorageCapacity;
      }

      return new StorageUsage(consumed, workspacesMaxStorageSum);
    }

    public async Task<UpdateUserResult> UpdateAsync(User subject, User newSubject, User user)
    {
      var response = new UpdateUserResult();
      
      newSubject.FullName = newSubject.FullName?.Trim();
      if (!string.IsNullOrWhiteSpace(newSubject.FullName) && newSubject.FullName != subject.FullName)
      {
        subject.FullName = newSubject.FullName;
      }
      else if (string.IsNullOrWhiteSpace(newSubject.FullName))
      {
        response.Errors.Add(Error.InvalidFullNameError);
      }

      newSubject.Email = newSubject.Email?.ToLowerInvariant().Trim();
      if (!string.IsNullOrWhiteSpace(newSubject.Email) && newSubject.Email != subject.Email)
      {
        List<Error> errors = await ValidateEmailAsync(subject.Email);
        if (errors.Any())
        {
          response.Errors.AddRange(errors);
        }
        else
        {
          subject.Email = newSubject.Email;
          if (_coreSettings.GravatarIntegration)
          {
            subject.Image = Base64ImageService.GenerateGravatar(newSubject.Email, _profileImageSize.Width);
          }
        }
      }
      else if (string.IsNullOrWhiteSpace(newSubject.Email))
      {
        response.Errors.Add(Error.InvalidEmailError);
      }

      // TODO: this is a bad idea, PasswordHash is used to represent plain text password and hashed one
      if (newSubject.PasswordHash != subject.PasswordHash &&
          _passwordService.HashPassword(newSubject.PasswordHash) != subject.PasswordHash &&
          !string.IsNullOrEmpty(newSubject.PasswordHash))
      {
        List<Error> errors = _passwordService.ValidatePassword(newSubject.PasswordHash).ToList();
        if (errors.Any())
        {
          response.Errors.AddRange(errors);
        }
        else
        {
          subject.PasswordHash = _passwordService.HashPassword(newSubject.PasswordHash);
        }
      }
      else if (string.IsNullOrEmpty(newSubject.PasswordHash))
      {
        response.Errors.Add(Error.NullOrEmptyPasswordError);
      }

      try
      {
        var proto = await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = user.Id,
          User = _mapper.Map<UserProto>(subject)
        });

        response.Result = _mapper.Map<User>(proto);

        return response;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> UpdateImageAsync(string imagePath, User subject, User user)
    {
      Image image = new Image(imagePath);

      if (image.Width > image.Height)
      {
        image.ScaleWithAspectRatio(_profileImageSize.Width, 0);
      }
      else
      {
        image.ScaleWithAspectRatio(0, _profileImageSize.Height);
      }

      using MemoryStream stream = new MemoryStream();
      image.SaveAsPngToStream(stream);
      subject.Image = Base64ImageService.Base64ImagePrefix + Convert.ToBase64String(stream.ToArray());

      try
      {
        var proto = await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = user.Id,
          User = _mapper.Map<UserProto>(subject)
        });

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task DeleteAsync(User subject, User user)
    {
      try
      {
        await _protoClient.DeleteAsync(new UserDeleteRequestProto
        {
          UserId = user.Id,
          User = _mapper.Map<UserProto>(subject)
        });
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> RefreshImageFromGravatarAsync(User subject, User user)
    {
      if (!_coreSettings.GravatarIntegration)
      {
        throw new Exception("Gravatar integration not enabled");
      }

      subject.Image = Base64ImageService.GenerateGravatar(subject.Email, _profileImageSize.Width);

      try
      {
        var proto = await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = user.Id,
          User = _mapper.Map<UserProto>(subject)
        });

        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindSystemUserAsync()
    {
      try
      {
        var proto = await _protoClient.FindSystemUserAsync(new UserFindSystemUserRequestProto());
        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<User> FindAdminUserAsync()
    {
      try
      {
        var proto = await _protoClient.FindAdminUserAsync(new UserFindAdminUserRequestProto());
        return _mapper.Map<User>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<List<Error>> ValidateEmailAsync(string email)
    {
      List<Error> errors = new List<Error>();

      try
      {
        User systemUser = await FindSystemUserAsync();
        await _protoClient.FindByFieldAsync(new UserFindByFieldRequestProto
        {
          UserId = systemUser.Id,
          Field = "email",
          Value = email
        });

        errors.Add(new Error(
          "AccountExists",
          "An account with this email already exists."));
      }
      catch
      {
        // ignored
      }

      /* For security reasons, we don't want to show that
         a given Email is both invalid and exists.
       */
      if (!EmailValidationService.IsValid(email))
      {
        return new List<Error> {Error.InvalidEmailError};
      }

      return errors;
    }

    public void CreateSystemUserIfNotExists()
    {
      try
      {
        FindSystemUserAsync().Wait();
      }
      catch (Exception)
      {
        string password = GenerateRandomPassword();
        System.IO.File.WriteAllText(Path.Combine(Directory.GetCurrentDirectory(), "config", "system_user_password.txt"), password);
        
        var user = new User
        {
          Username = "system",
          Email = "system@localhost",
          FullName = "System",
          IsSystem = true,
          IsSuperuser = true,
          PasswordHash = password,
          IsEmailConfirmed = true
        };
        InsertInsecureAsync(user).Wait();
      }
    }

    public void CreateAdminUserIfNotExists()
    {
      try
      {
        FindAdminUserAsync().Wait();
      }
      catch (Exception)
      {
        string password = GenerateRandomPassword();
        System.IO.File.WriteAllText(Path.Combine(Directory.GetCurrentDirectory(), "config", "admin_user_password.txt"), password);
        
        var user = new User
        {
          Username = "admin",
          Email = "admin@localhost",
          FullName = "Admin",
          IsSuperuser = true,
          PasswordHash = password,
          IsEmailConfirmed = true
        };
        InsertInsecureAsync(user).Wait();
      }
    }

    /// <summary>
    /// Generates a Random Password
    /// respecting the given strength requirements.
    /// </summary>
    /// <param name="opts">A valid PasswordOptions object
    /// containing the password strength requirements.</param>
    /// <returns>A random password</returns>
    private string GenerateRandomPassword(PasswordOptions opts = null)
    {
      if (opts == null)
        opts = new PasswordOptions
        {
          RequiredLength = 12,
          RequiredUniqueChars = 12,
          RequireDigit = true,
          RequireLowercase = true,
          RequireNonAlphanumeric = true,
          RequireUppercase = true
        };

      string[] randomChars =
      {
        "ABCDEFGHJKLMNOPQRSTUVWXYZ", // uppercase 
        "abcdefghijkmnopqrstuvwxyz", // lowercase
        "0123456789", // digits
        "!@$?_-#%&^" // non-alphanumeric
      };

      Random rand = new Random(Environment.TickCount);
      List<char> chars = new List<char>();

      if (opts.RequireUppercase)
      {
        chars.Insert(rand.Next(0, chars.Count), randomChars[0][rand.Next(0, randomChars[0].Length)]);
      }

      if (opts.RequireLowercase)
      {
        chars.Insert(rand.Next(0, chars.Count), randomChars[1][rand.Next(0, randomChars[1].Length)]);
      }

      if (opts.RequireDigit)
      {
        chars.Insert(rand.Next(0, chars.Count), randomChars[2][rand.Next(0, randomChars[2].Length)]);
      }

      if (opts.RequireNonAlphanumeric)
      {
        chars.Insert(rand.Next(0, chars.Count), randomChars[3][rand.Next(0, randomChars[3].Length)]);
      }

      for (int i = chars.Count; i < opts.RequiredLength || chars.Distinct().Count() < opts.RequiredUniqueChars; i++)
      {
        string rcs = randomChars[rand.Next(0, randomChars.Length)];
        chars.Insert(rand.Next(0, chars.Count), rcs[rand.Next(0, rcs.Length)]);
      }

      return new string(chars.ToArray());
    }
  }
}