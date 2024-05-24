namespace Defyle.Core.Storage.Services
{
  using System.IO;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Exceptions;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Inode.Models;
  using Inode.Services;
  using Ocr.Queue;
  using Preview.Queue;
  using S3.Services;
  using Workspace.Models;
  using File = Models.File;

  public class StorageService : QueueService
	{
		private readonly CoreSettings _coreSettings;
    private readonly PathService _pathService;
    private readonly InodeEngine _inodeEngine;
    private readonly InodeNotificationService _inodeNotificationService;
    private readonly PreviewQueueService _previewQueueService;
    private readonly OcrQueueService _ocrQueueService;
    private readonly FileService _fileService;

    public StorageService(
			CoreSettings coreSettings,
      PathService pathService,
      InodeEngine inodeEngine,
      InodeNotificationService inodeNotificationService,
      PreviewQueueService previewQueueService,
      OcrQueueService ocrQueueService,
      FileService fileService)
			: base(coreSettings.MessageBroker)
		{
			_coreSettings = coreSettings;
      _pathService = pathService;
      _inodeEngine = inodeEngine;
      _inodeNotificationService = inodeNotificationService;
      _previewQueueService = previewQueueService;
      _ocrQueueService = ocrQueueService;
      _fileService = fileService;
    }

		public async Task StoreAsync(Workspace workspace, InodeFacet inode, File file,
      Stream stream, string cipherPassword, User user)
		{
			string storageDirectory = _pathService.GetFileDirectory(workspace, file);

			if (!Directory.Exists(storageDirectory))
			{
        Directory.CreateDirectory(storageDirectory);
			}

      await _inodeEngine.SetFileAsync(inode, file, user);

			string path = _pathService.GetOriginalFile(workspace, file);

			if (workspace.Encrypted)
			{
				if (string.IsNullOrWhiteSpace(cipherPassword))
				{
					throw new InternalServerErrorException().WithError(Error.WorkspacePasswordRequiredError);
				}

				byte[] key = EncryptionService.EncryptionKey(workspace, cipherPassword);

				// This should be placed after _nodeTree.SetFile(),
				// otherwise the cleanup service will end up deleting this physical file
        using FileStream fileStream = new FileStream(path, FileMode.Create);
        EncryptionService.EncryptStreamWithSalt(stream, fileStream, key);
      }
			else
			{
				// This should be placed after _nodeTree.SetFile(),
				// otherwise the cleanup service will end up deleting this physical file
				using (FileStream fileStream = new FileStream(path, FileMode.Create))
				{
					await stream.CopyToAsync(fileStream);
				}

				Inode updated = await _inodeEngine.FindByIdAsync(inode.Id, user);
				await QueueFilePreviewAsync(workspace, updated, user);
      }

      if (_coreSettings.S3Enabled)
      {
        S3Service s3Service = new S3Service(_coreSettings.S3);
        await s3Service.UploadAsync(path, _pathService.GetS3OriginalFile(workspace, file));
      }
    }

		private async Task QueueFilePreviewAsync(Workspace workspace, Inode inode, User user)
    {
      File file = await _fileService.FindLatestForInodeOrNullAsync(inode, user);
      if (file == null)
      {
        return;
      }
      
      if (file.Size > _coreSettings.PreviewFileSizeLimit)
      {
        return;
      }
      
      string fileType = await _fileService.GetFileTypeAsync(file.Mime);
      string fileCategory = await _fileService.GetFileCategoryAsync(fileType);

      if (fileCategory == "image")
			{
        await _fileService.SetPropertyAsync(file, "file.imagePreview.status", "pending", user);
        await _previewQueueService.SendImagePreviewMessageAsync(workspace, inode, user);

        if (_coreSettings.EnableTileMapPreview)
        {
          await _fileService.SetPropertyAsync(file, "file.tileMap.status", "pending", user);
          await _previewQueueService.SendTileMapMessageAsync(workspace, inode, user);
        }
        
        await _inodeNotificationService.SendInodesPropertiesUpdatedAsync(workspace.Id, new[] {inode.Id});
      }
			else if (fileCategory == "document")
			{
				if (file.GetMime() == "application/pdf")
				{
					await _fileService.SetPropertyAsync(file, "file.documentPreview.status", "ready", user);

          if (file.IndexContent)
          {
            await _fileService.SetPropertyAsync(file, "file.ocr.searchablePdf.status", "pending", user);
            await _ocrQueueService.SendSearchablePdfMessageAsync(workspace, inode, user);
          }
          else
          {
            await _fileService.SetPropertyAsync(file, "file.ocr.textExtraction.status", "pending", user);
            await _ocrQueueService.SendTextExtractionMessageAsync(workspace, inode, user);
          }
        }
				else
				{
          await _fileService.SetPropertyAsync(file, "file.documentPreview.status", "pending", user);
          await _previewQueueService.SendDocumentMessageAsync(workspace, inode, user);
        }
        
        await _inodeNotificationService.SendInodesPropertiesUpdatedAsync(workspace.Id, new[] {inode.Id});
      }
		}
	}
}