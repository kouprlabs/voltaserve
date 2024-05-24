namespace Defyle.Core.Inode.Services
{
  using System;
  using System.Collections.Generic;
  using System.IO;
  using System.Linq;
  using System.Threading.Tasks;
  using Auth.Models;
  using Auth.Services;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Microsoft.AspNetCore.Http;
  using Models;
  using Pocos;
  using Storage.Exceptions;
  using Storage.Models;
  using Storage.Services;
  using Workspace.Models;
  using File = Storage.Models.File;

  public class InodeService : QueueService
	{
		private readonly CoreSettings _coreSettings;
    private readonly PathService _pathService;
		private readonly StorageService _storageService;
    private readonly InodeEngine _inodeEngine;
    private readonly UserService _userService;
    private readonly InodeNotificationService _inodeNotificationService;
    private readonly FileService _fileService;

    public InodeService(
			CoreSettings coreSettings,
      PathService pathService,
			StorageService storageService,
      InodeEngine inodeEngine,
      UserService userService,
      InodeNotificationService inodeNotificationService,
      FileService fileService)
			: base(coreSettings.MessageBroker)
		{
			_coreSettings = coreSettings;
      _pathService = pathService;
			_storageService = storageService;
      _inodeEngine = inodeEngine;
      _userService = userService;
      _inodeNotificationService = inodeNotificationService;
      _fileService = fileService;
    }

		public async Task<InodeFacet> CreateFileAsync(Workspace workspace, string parentId, string name, User user)
    {
      InodeFacet parent = null;
      if (!string.IsNullOrWhiteSpace(parentId))
      {
        parent = await _inodeEngine.FindByIdAsync(parentId, user);
      }

      InodeFacet inode = await _inodeEngine.CreateFileAsync(workspace, parent, name, user);

      if (parent != null)
      {
        await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, new[] {parent.Id});
      }

      return inode;
		}
    
    public async Task<Inode> CreateDirectoryAsync(Workspace workspace, string parentId, string name, User user)
    {
      InodeFacet parent = null;
      if (!string.IsNullOrWhiteSpace(parentId))
      {
        parent = await _inodeEngine.FindByIdAsync(parentId, user);
      }

      Inode inode = await _inodeEngine.CreateDirectoryAsync(workspace, parent, name, user);

      if (parent != null)
      {
        await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, new[] {parent.Id});
      }

      return inode;
    }

		public async Task<InodeFacet> CreateFromFileAsync(Workspace workspace, InodeFacet parent,
      bool requireOcr, string cipherPassword, IFormFile formFile, User user)
		{
			await FitsAllowedStorageUsageOrThrowAsync(formFile.Length, workspace);

      InodeFacet inode = await _inodeEngine.CreateFileAsync(workspace, parent, formFile.FileName, user);

			string extension = Path.GetExtension(formFile.FileName);
			if (string.IsNullOrWhiteSpace(extension))
			{
				extension = null;
			}
			else
			{
				extension = extension.Substring(1);
			}

			File file = new File(extension, formFile.Length, requireOcr);

			await _storageService.StoreAsync(workspace, inode, file, formFile.OpenReadStream(), cipherPassword, user);
      
      await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, new[] {parent.Id});

			return inode;
		}

		public async Task<InodeFacet> FindOneAsync(string inodeId, User user)
    {
      return await _inodeEngine.FindByIdAsync(inodeId, user);
		}
    
    public async Task<InodeFacet> FindRootAsync(Workspace workspace, User user)
    {
      return await _inodeEngine.FindRootAsync(workspace, user);
    }

		public async Task<DeleteInodesResult> DeleteManyAsync(
      Workspace workspace, IEnumerable<string> ids, User user)
		{
      var response = new DeleteInodesResult();

      List<string> physicalFilePathsToDelete = new List<string>();
			foreach (string id in ids)
			{
        try
        {
          InodeFacet inode = await _inodeEngine.FindByIdAsync(id, user);
          IEnumerable<string> deletedFileIds = await _inodeEngine.DeleteAsync(inode, user);
          
          IEnumerable<string> paths = deletedFileIds.Select(path => Path.Combine(Path.Combine(_coreSettings.DataDirectory, workspace.Id), "file-data", path));
          physicalFilePathsToDelete.AddRange(paths);
          
          response.Succeeded.Add(inode.Id);
          response.AffectedParentIds.Add(inode.ParentId);
        }
        catch (Exception)
        {
          response.Failed.Add(id);
          response.Errors.Add(Error.InodeNotFoundError
            .WithDescription($"an error occured when processing '{id}'"));
        }
      }
      
			if (physicalFilePathsToDelete.Count > 0 && !_coreSettings.S3Enabled)
			{
        var message = new CleanupQueueMessage
        {
          Id = Guid.NewGuid().ToString(),
          Paths = physicalFilePathsToDelete 
        };
        SendQueueMessage(message, _coreSettings.CleanupQueue);
			}
      
      await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, response.AffectedParentIds);
      
      response.AffectedParentIds = response.AffectedParentIds.Distinct().ToList();
      response.Succeeded = response.Succeeded.Distinct().ToList();
      response.Failed = response.Failed.Distinct().ToList();

			return response;
		}

		public async Task<MoveInodesResult> MoveManyAsync(Workspace workspace, InodeFacet inode,
      IEnumerable<string> ids, User user)
		{
      if (inode.Type != Inode.InodeTypeDirectory)
			{
				throw new InternalServerErrorException().WithError(Error.DirectoryOnlyOperationError);
			}

			var response = new MoveInodesResult();

      foreach (string id in ids)
			{
        try
        {
          InodeFacet source = await _inodeEngine.FindByIdAsync(id, user);
          await _inodeEngine.MoveAsync(source, inode, user);
          
          Inode updatedSource = await _inodeEngine.FindByIdAsync(id, user);

          response.Succeeded.Add(source.Id);
          response.PreviousParentIds.Add(source.ParentId);
          response.CurrentParentIds.Add(updatedSource.ParentId);
        }
        catch (Exception)
        {
          response.Failed.Add(id);
          response.Errors.Add(Error.InodeNotFoundError
            .WithDescription($"an error occured when processing '{id}'"));
        }
			}
      
			List<string> updatedParents = new List<string>();
			updatedParents.AddRange(response.PreviousParentIds);
			updatedParents.AddRange(response.CurrentParentIds);
      
      await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, updatedParents);

      response.CurrentParentIds = response.CurrentParentIds.Distinct().ToList();
      response.PreviousParentIds = response.PreviousParentIds.Distinct().ToList();
      response.Succeeded = response.Succeeded.Distinct().ToList();
      response.Failed = response.Failed.Distinct().ToList();

			return response;
		}
    
    public async Task<MoveInodesChildrenResult> MoveManyChildrenAsync(Workspace workspace, InodeFacet destination,
      IEnumerable<string> ids, User user)
    {
      if (destination.Type != Inode.InodeTypeDirectory)
      {
        throw new InternalServerErrorException().WithError(
          Error.DirectoryOnlyOperationError.WithDescription($"Destination '{destination.Id}' is not a directory."));
      }

      var response = new MoveInodesChildrenResult();

      foreach (string id in ids)
      {
        InodeFacet source = await _inodeEngine.FindByIdAsync(id, user);
        if (source.Type != Inode.InodeTypeDirectory)
        {
          response.Failed.Add(source.Id);
          response.Errors.Add(Error.DirectoryOnlyOperationError
            .WithDescription($"Inode '{source.Id}' is not a directory."));
          continue;
        }

        try
        {
          await _inodeEngine.MoveChildrenAsync(source, destination, user);
          response.Succeeded.Add(source.Id);
          response.PreviousParentIds.Add(source.Id);
        }
        catch (Exception)
        {
          response.Failed.Add(id);
          response.Errors.Add(Error.InodeNotFoundError
            .WithDescription($"an error occured when processing '{id}'"));
        }
      }
      
      response.CurrentParentId = destination.Id;

      List<string> updatedParents = new List<string>();
      updatedParents.AddRange(response.PreviousParentIds);
      updatedParents.Add(response.CurrentParentId);
      
      await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, updatedParents);
      
      response.PreviousParentIds = response.PreviousParentIds.Distinct().ToList();
      response.Succeeded = response.Succeeded.Distinct().ToList();
      response.Failed = response.Failed.Distinct().ToList();

      return response;
    }

		public async Task<CopyInodesResult> CopyManyAsync(Workspace workspace,
      InodeFacet inode, IEnumerable<string> ids, User user)
		{
      if (inode.Type != Inode.InodeTypeDirectory)
			{
				throw new InternalServerErrorException().WithError(Error.DirectoryOnlyOperationError);
			}

			var response = new CopyInodesResult();

      foreach (string id in ids)
			{
				try
				{
          InodeFacet subject = await _inodeEngine.FindByIdAsync(id, user);
          string createdId = await _inodeEngine.CopyAsync(subject, inode, user);
          response.Created.Add(createdId);
          response.Succeeded.Add(id);
				}
				catch (Exception)
				{
          response.Errors.Add(Error.InodeNotFoundError.WithDescription($"an error occured when processing '{id}'"));
          response.Failed.Add(id);
				}
			}

      await _inodeNotificationService.SendInodesChildrenUpdatedAsync(workspace.Id, new List<string> {inode.Id});

      response.Created = response.Created.Distinct().ToList();
      response.Succeeded = response.Succeeded.Distinct().ToList();
      response.Failed = response.Failed.Distinct().ToList();
      
			return response;
		}

		public async Task<InodeFacet> UpdateNameAsync(Workspace workspace, InodeFacet inode, string name, User user)
		{
      await _inodeEngine.SetNameAsync(inode, name, user);
      
      await _inodeNotificationService.SendInodesPropertiesUpdatedAsync(workspace.Id, new[] {inode.Id});

			return inode;
		}

		public async Task<InodePagedResult> GetChildrenAsync(InodeFacet inode, int page, int count, User user)
		{
      if (inode.Type != Inode.InodeTypeDirectory)
			{
				throw new InternalServerErrorException().WithError(Error.DirectoryOnlyOperationError);
			}

      InodePagedResult pagedResult;
			if (user.IsSuperuser)
			{
				pagedResult = await _inodeEngine.GetChildrenAsync(inode, page, count, user);
			}
			else
			{
				pagedResult = await _inodeEngine.GetChildrenAsync(inode, page, count, user);
			}

      return pagedResult;
		}

		public async Task<InodePagedResult> SearchAsync(Workspace workspace, InodeSearchOptions request, int page, int count, User user)
    {
      InodePagedResult pagedResult = await _inodeEngine.SearchAsync(workspace, page, count, request, user);

			return pagedResult;
		}
    
    public async Task<long> CountChildrenAsync(InodeFacet inode, User user)
    {
      if (inode.Type != Inode.InodeTypeDirectory)
      {
        throw new InternalServerErrorException().WithError(Error.DirectoryOnlyOperationError);
      }
      
      return await _inodeEngine.CountChildrenAsync(inode, user);
    }

    public async Task UpdateFileAsync(Workspace workspace, InodeFacet inode, string filename,
      long length, bool indexContent, Stream inputStream, string cipherPassword, User user)
		{
			await FitsAllowedStorageUsageOrThrowAsync(length, workspace);
      
			string extension = Path.GetExtension(filename).Substring(1);

			File file = new File(extension, length, indexContent);

			await _storageService.StoreAsync(workspace, inode, file, inputStream, cipherPassword, user);
      
      await _inodeNotificationService.SendInodesPropertiesUpdatedAsync(workspace.Id, new[] {inode.Id});
    }

    /// <returns>Storage usage in bytes</returns>
		public async Task<StorageUsage> GetStorageUsageAsync(Workspace workspace, InodeFacet inode, User user)
    {
      long consumed = await _inodeEngine.GetStorageUsageAsync(inode, user);
      return new StorageUsage(consumed, workspace.StorageCapacity);;
		}

		private async Task FitsAllowedStorageUsageOrThrowAsync(long length, Workspace workspace)
    {
      var systemUser = await _userService.FindSystemUserAsync();
      var rootNode = await _inodeEngine.FindRootAsync(workspace, systemUser);
      
      long consumed = await _inodeEngine.GetStorageUsageAsync(rootNode, systemUser);
      long expectedConsumption = consumed + length;
      if (expectedConsumption > workspace.StorageCapacity)
			{
				throw new StorageUsageExceededException().WithError(Error.StorageUsageExceededError);
			}
		}
  }
}