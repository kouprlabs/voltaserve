namespace Defyle.Core.Workspace.Services
{
  using System;
  using System.Collections.Generic;
  using System.Drawing;
  using System.IO;
  using System.Linq;
  using System.Security.Cryptography;
  using System.Text;
  using System.Threading.Tasks;
  using Auth.Models;
  using Auth.Services;
  using AutoMapper;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Inode.Models;
  using Inode.Services;
  using Microsoft.Extensions.Logging;
  using Pocos;
  using Preview.Models;
  using Preview.Services;
  using Proto;
  using Storage.Services;
  using Workspace = Models.Workspace;

  public class WorkspaceService
	{
    private readonly Size _imageSize = new Size(230, 230);
		private readonly CoreSettings _coreSettings;
    private readonly InodeEngine _inodeEngine;
    private readonly PasswordService _passwordService;
    private readonly WorkspaceNotificationService _notificationService;
    private readonly WorkspaceServiceProto.WorkspaceServiceProtoClient _protoClient;
    private readonly GrpcExceptionTranslator _exceptionTranslator;
    private readonly IMapper _mapper;
    private readonly ILogger<WorkspaceService> _logger;

    public WorkspaceService(
			CoreSettings coreSettings,
      InodeEngine inodeEngine,
      PasswordService passwordService,
      WorkspaceNotificationService notificationService,
      WorkspaceServiceProto.WorkspaceServiceProtoClient protoClient,
      GrpcExceptionTranslator exceptionTranslator,
      IMapper mapper,
      ILogger<WorkspaceService> logger)
    {
			_coreSettings = coreSettings;
      _inodeEngine = inodeEngine;
      _passwordService = passwordService;
      _notificationService = notificationService;
      _protoClient = protoClient;
      _exceptionTranslator = exceptionTranslator;
      _mapper = mapper;
      _logger = logger;
    }

    public async Task<Workspace> InsertAsync(CreateWorkspaceOptions options, User user)
    {
      var insertOptions = new WorkspaceInsertRequestProto
      {
        UserId = user.Id,
        Workspace = new WorkspaceProto
        {
          PartitionId = _coreSettings.PartitionId,
          Name = options.Name,
          Encrypted = options.Encrypted,
          StorageCapacity = _coreSettings.DefaultWorkspaceStorageCapacity
        }
      };

      if (!string.IsNullOrWhiteSpace(options.Image))
      {
        insertOptions.Workspace.Image = options.Image;
      }

      if (options.Encrypted)
      {
        List<Error> errors = _passwordService.ValidatePassword(options.Password);

        if (errors.Any())
        {
          throw new InternalServerErrorException().WithErrors(errors);
        }

        // Generate random salt
        byte[] salt = new byte[16];
        new RNGCryptoServiceProvider().GetBytes(salt);

        // Generate random encryption key
        byte[] key = new byte[16];
        new RNGCryptoServiceProvider().GetBytes(key);

        // Encrypt key with user password
        byte[] cipherKey = EncryptionService.EncryptBytesWithSalt(
          key,
          Encoding.ASCII.GetBytes(options.Password),
          salt);

        // Generate random transit key
        byte[] transitKey = new byte[16];
        new RNGCryptoServiceProvider().GetBytes(transitKey);

        // Generate random transit IV
        byte[] transitIv = new byte[16];
        new RNGCryptoServiceProvider().GetBytes(transitIv);

        insertOptions.Workspace.Encrypted = true;
        insertOptions.Workspace.PasswordHash = _passwordService.HashPassword(options.Password);
        insertOptions.Workspace.CipherKey = Convert.ToBase64String(cipherKey);
        insertOptions.Workspace.Salt = Convert.ToBase64String(salt);
        insertOptions.Workspace.TransitKey = Convert.ToBase64String(transitKey);
        insertOptions.Workspace.TransitIv = Convert.ToBase64String(transitIv);
      }

      try
      {
        var proto = await _protoClient.InsertAsync(insertOptions);
        Workspace workspace = _mapper.Map<Workspace>(proto);
        
        await _notificationService.SendWorkspacesUpdatedAsync();

        return workspace;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

    public async Task<Workspace> FindAsync(string id, User user)
    {
      try
      {
        var request = new WorkspaceFindRequestProto {UserId = user.Id, Id = id};
        WorkspaceProto proto = await _protoClient.FindAsync(request);
      
        return _mapper.Map<Workspace>(proto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
    
    public async Task<IEnumerable<Workspace>> FindAllAsync(User user)
    {
      try
      {
        var request = new WorkspaceFindAllRequestProto {UserId = user.Id};
        WorkspaceFindAllResponseProto response = await _protoClient.FindAllAsync(request);

        return response.Workspaces.Select(p => _mapper.Map<Workspace>(p)).ToList();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

		public async Task<Workspace> UpdateNameAsync(Workspace workspace, string name, User user)
    {
      try
      {
        workspace.Name = name;
        
        var updateRequestProto = new WorkspaceUpdateRequestProto {UserId = user.Id, Workspace = _mapper.Map<WorkspaceProto>(workspace)};
        var proto = await _protoClient.UpdateAsync(updateRequestProto);
        var result = _mapper.Map<Workspace>(proto);

        await _notificationService.SendWorkspacesUpdatedAsync();
      
        return result;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
		}

		public async Task<Workspace> UpdateImageAsync(Workspace workspace, string path, User user)
		{
      Image image = new Image(path);

			if (image.Width > image.Height)
			{
				image.ScaleWithAspectRatio(_imageSize.Width, 0);
			}
			else
			{
				image.ScaleWithAspectRatio(0, _imageSize.Height);
			}

      await using MemoryStream stream = new MemoryStream();
      image.SaveAsPngToStream(stream);
      string base64 = Base64ImageService.Base64ImagePrefix + Convert.ToBase64String(stream.ToArray());

      workspace.Image = base64;

      try
      {
        var updateRequestProto = new WorkspaceUpdateRequestProto {UserId = user.Id, Workspace = _mapper.Map<WorkspaceProto>(workspace)};
        var proto = await _protoClient.UpdateAsync(updateRequestProto);
        var result = _mapper.Map<Workspace>(proto);
        
        await _notificationService.SendWorkspacesUpdatedAsync();
        
        return result;
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

		public async Task UpdatePasswordAsync(Workspace workspace, string currentPassword, string newPassword, User user)
		{
			if (!workspace.Encrypted)
			{
				throw new InternalServerErrorException().WithError(Error.WorkspaceNotEncryptedError);
			}
      
			if (!_passwordService.VerifyHashedPassword(workspace.PasswordHash, currentPassword))
			{
				throw new InternalServerErrorException().WithError(Error.InvalidWorkspacePasswordError);
			}

			List<Error> errors = _passwordService.ValidatePassword(newPassword);

			if (errors.Any())
			{
				throw new InternalServerErrorException().WithErrors(errors);
			}

			// Get encryption key using current password
			byte[] key = EncryptionService.DecryptBytesWithSalt(
				Convert.FromBase64String(workspace.CipherKey),
				Encoding.UTF8.GetBytes(currentPassword),
				Convert.FromBase64String(workspace.Salt));

			// Encrypt key with the new user password
			byte[] newCipherKey = EncryptionService.EncryptBytesWithSalt(
				key,
				Encoding.ASCII.GetBytes(newPassword),
				Convert.FromBase64String(workspace.Salt));

			string newPasswordHash = _passwordService.HashPassword(newPassword);

      workspace.CipherKey = Convert.ToBase64String(newCipherKey);
      workspace.PasswordHash = newPasswordHash;

      try
      {
        var updateRequestProto = new WorkspaceUpdateRequestProto {UserId = user.Id, Workspace = _mapper.Map<WorkspaceProto>(workspace)};
        await _protoClient.UpdateAsync(updateRequestProto);
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }

		public void VerifyPassword(Workspace workspace, string cipherPassword)
		{
			if (!workspace.Encrypted)
			{
				throw new InternalServerErrorException().WithError(Error.WorkspaceNotEncryptedError);
			}

      string password = EncryptionService.DecryptBytes(
				Convert.FromBase64String(cipherPassword),
				Convert.FromBase64String(workspace.TransitKey),
				Convert.FromBase64String(workspace.TransitIv));

			if (!_passwordService.VerifyHashedPassword(workspace.PasswordHash, password))
			{
				throw new InternalServerErrorException().WithError(Error.InvalidWorkspacePasswordError);
			}
		}

		public async Task DeleteAsync(Workspace workspace, User user)
		{
      try
      {
        InodeFacet rootInode = await _inodeEngine.FindRootAsync(workspace, user);
        await _inodeEngine.DeleteAsync(rootInode, user);
      }
      catch (ResourceNotFoundException e)
      {
        _logger.LogCritical(e, $"Workspace {workspace.Id} has no root inode");
      }

      try
      {
        var request = new WorkspaceDeleteRequestProto {UserId = user.Id, Workspace = _mapper.Map<WorkspaceProto>(workspace)};
        await _protoClient.DeleteAsync(request);

        await _notificationService.SendWorkspacesUpdatedAsync();
      }
      catch (Exception e)
      {
        throw _exceptionTranslator.Translate(e);
      }
    }
  }
}