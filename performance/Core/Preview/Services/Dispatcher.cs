namespace Defyle.Core.Preview.Services
{
  using System;
  using System.Text;
  using System.Threading;
  using Auth.Models;
  using Auth.Services;
  using Infrastructure.Poco;
  using Infrastructure.Services;
  using Inode.Models;
  using Inode.Pocos;
  using Inode.Services;
  using Microsoft.Extensions.Logging;
  using Newtonsoft.Json;
  using Ocr.Queue;
  using Ocr.Services;
  using Queue;
  using RabbitMQ.Client;
  using RabbitMQ.Client.Events;
  using Storage.Models;
  using Storage.Services;
  using Workspace.Models;
  using Workspace.Services;

  public class Dispatcher : QueueService
	{
		private readonly CoreSettings _coreSettings;
    private readonly InodeEngine _inodeEngine;
    private readonly TextExtractionService _textExtractionService;
    private readonly InodeNotificationService _inodeNotificationService;
    private readonly OcrQueueService _ocrQueueService;
    private readonly UserService _userService;
    private readonly WorkspaceService _workspaceService;
    private readonly FileService _fileService;
    private readonly PathService _pathService;
    private readonly ILogger<Dispatcher> _logger;
    private ConnectionFactory _factory;
		private IConnection _connection;
		private IModel _imageStatusModel;
		private IModel _tileMapStatusModel;
		private IModel _documentStatusModel;
    private IModel _textExtractionStatusModel;
    private IModel _searchablePdfStatusModel;

    public Dispatcher(
			CoreSettings coreSettings,
      InodeEngine inodeEngine,
      TextExtractionService textExtractionService,
      InodeNotificationService inodeNotificationService,
      OcrQueueService ocrQueueService,
      UserService userService,
      WorkspaceService workspaceService,
      FileService fileService,
      PathService pathService,
      ILogger<Dispatcher> logger) : base(coreSettings.MessageBroker)
		{
			_coreSettings = coreSettings;
      _inodeEngine = inodeEngine;
      _textExtractionService = textExtractionService;
      _inodeNotificationService = inodeNotificationService;
      _ocrQueueService = ocrQueueService;
      _userService = userService;
      _workspaceService = workspaceService;
      _fileService = fileService;
      _pathService = pathService;
      _logger = logger;
    }

		public void Register()
    {
      _factory = new ConnectionFactory
      {
        Uri = new Uri(_coreSettings.MessageBroker.Url)
      };

      while (true)
      {
        try
        {
          _connection = _factory.CreateConnection();
          _logger.LogInformation("RabbitMQ connected.");
          break;
        }
        catch (RabbitMQ.Client.Exceptions.BrokerUnreachableException e)
        {
          _logger.LogCritical(e, $"Connection to RabbitMQ failed with error: {e.Message}. Retrying in 3 seconds...");
          Thread.Sleep(3000);
        }
      }

      _imageStatusModel = SetupQueue(OnImageStatus, _coreSettings.PreviewWorker.ImageStatusQueue);
			_tileMapStatusModel = SetupQueue(OnTileMapStatus, _coreSettings.PreviewWorker.TileMapStatusQueue);
			_documentStatusModel = SetupQueue(OnDocumentStatus, _coreSettings.PreviewWorker.DocumentStatusQueue);
      _textExtractionStatusModel = SetupQueue(OnTextExtractionStatus, _coreSettings.OcrWorker.TextExtractionStatusQueue);
      _searchablePdfStatusModel = SetupQueue(OnSearchablePdfStatus, _coreSettings.OcrWorker.SearchablePdfStatusQueue);
    }

		public void Deregister()
		{
			_imageStatusModel.Close();
			_tileMapStatusModel.Close();
			_documentStatusModel.Close();
      _textExtractionStatusModel.Close();
      _searchablePdfStatusModel.Close();
			_connection.Close();
		}

		private IModel SetupQueue(EventHandler<BasicDeliverEventArgs> eventHandler, string queue)
		{
			IModel model = _connection.CreateModel();

			model.QueueDeclare(queue, true, false, false, null);

			model.BasicQos(0, 1, false);

			EventingBasicConsumer consumer = new EventingBasicConsumer(model);
			consumer.Received += eventHandler;

			model.BasicConsume(queue, false, consumer);

			return model;
		}

		private void OnImageStatus(object model, BasicDeliverEventArgs args)
		{
      string rawMessage = Encoding.UTF8.GetString(args.Body);
			PreviewStatusMessage message = JsonConvert.DeserializeObject<PreviewStatusMessage>(rawMessage);
			string indentedMessage = JsonConvert.SerializeObject(message, Formatting.Indented);
			_logger.LogTrace($"Received on queue '{_coreSettings.PreviewWorker.ImageStatusQueue}': {indentedMessage}");
      
			try
			{
        User systemUser = _userService.FindSystemUserAsync().Result;
        Workspace workspace = _workspaceService.FindAsync(message.WorkspaceId, systemUser).Result;
        Inode inode = _inodeEngine.FindByIdAsync(message.InodeId, systemUser).Result;
        
        var file = _fileService.FindLatestForInodeOrNullAsync(inode, systemUser).Result;
        _fileService.SetPropertyAsync(file, "file.imagePreview.status", message.Status, systemUser).Wait();

        if (message.Status == "ready" && file != null && file.IndexContent)
        {
          _ocrQueueService.SendSearchablePdfMessageAsync(workspace, inode, systemUser).Wait();
          _fileService.SetPropertyAsync(file, "file.ocr.searchablePdf.status", "pending", systemUser).Wait();
        }
        
        _inodeNotificationService.SendInodesPropertiesUpdatedAsync(message.WorkspaceId, new[] {message.InodeId}).Wait();

				_logger.LogTrace($"Succeeded on queue '{_coreSettings.PreviewWorker.ImageStatusQueue}': {indentedMessage}");

				_imageStatusModel.BasicAck(args.DeliveryTag, false);
        
        if (_coreSettings.S3Enabled)
        {
          DeleteOriginalFileIfNoPendingStatus(workspace, file, systemUser);
        }
			}
      catch (Exception e)
			{
        _logger.LogCritical(e, $"Failed on queue '{_coreSettings.PreviewWorker.ImageStatusQueue}': {indentedMessage}");
        _imageStatusModel.BasicNack(args.DeliveryTag, false, false);
      }
		}

		private void OnTileMapStatus(object model, BasicDeliverEventArgs args)
		{
			string rawMessage = Encoding.UTF8.GetString(args.Body);
			PreviewStatusMessage message = JsonConvert.DeserializeObject<PreviewStatusMessage>(rawMessage);
			string indentedMessage = JsonConvert.SerializeObject(message, Formatting.Indented);
			_logger.LogTrace($"Received on queue '{_coreSettings.PreviewWorker.TileMapStatusQueue}': {indentedMessage}");

			try
			{
        User systemUser = _userService.FindSystemUserAsync().Result;
        Workspace workspace = _workspaceService.FindAsync(message.WorkspaceId, systemUser).Result;
        Inode inode = _inodeEngine.FindByIdAsync(message.InodeId, systemUser).Result;

        var file = _fileService.FindLatestForInodeOrNullAsync(inode, systemUser).Result;
        _fileService.SetPropertyAsync(file, "file.tileMap.status", message.Status, systemUser).Wait();
        
        _inodeNotificationService.SendInodesPropertiesUpdatedAsync(message.WorkspaceId, new[] {message.InodeId}).Wait();

				_logger.LogTrace($"Succeeded on queue '{_coreSettings.PreviewWorker.TileMapStatusQueue}': {indentedMessage}");

				_tileMapStatusModel.BasicAck(args.DeliveryTag, false);
        
        if (_coreSettings.S3Enabled)
        {
          DeleteOriginalFileIfNoPendingStatus(workspace, file, systemUser);
        }
			}
      catch (Exception e)
			{
        _logger.LogCritical(e, $"Failed on queue '{_coreSettings.PreviewWorker.TileMapStatusQueue}': {indentedMessage}");
        _tileMapStatusModel.BasicNack(args.DeliveryTag, false, false);
      }
		}

		private void OnDocumentStatus(object model, BasicDeliverEventArgs args)
		{
			string rawMessage = Encoding.UTF8.GetString(args.Body);
			PreviewStatusMessage message = JsonConvert.DeserializeObject<PreviewStatusMessage>(rawMessage);
			string indentedMessage = JsonConvert.SerializeObject(message, Formatting.Indented);
			_logger.LogTrace($"Received on queue '{_coreSettings.PreviewWorker.DocumentStatusQueue}': {indentedMessage}");

			try
			{
        User systemUser = _userService.FindSystemUserAsync().Result;
        Workspace workspace = _workspaceService.FindAsync(message.WorkspaceId, systemUser).Result;
        Inode inode = _inodeEngine.FindByIdAsync(message.InodeId, systemUser).Result;

        var file = _fileService.FindLatestForInodeOrNullAsync(inode, systemUser).Result;
        _fileService.SetPropertyAsync(file, "file.documentPreview.status", message.Status, systemUser).Wait();

        if (message.Status == "ready")
        {
          _ocrQueueService.SendTextExtractionMessageAsync(workspace, inode, systemUser).Wait();
          _fileService.SetPropertyAsync(file, "file.ocr.textExtraction.status", "pending", systemUser).Wait();
        }
        
        _inodeNotificationService.SendInodesPropertiesUpdatedAsync(message.WorkspaceId, new[] {message.InodeId}).Wait();

				_logger.LogTrace($"Succeeded on queue '{_coreSettings.PreviewWorker.DocumentStatusQueue}': {indentedMessage}");

				_documentStatusModel.BasicAck(args.DeliveryTag, false);
        
        if (_coreSettings.S3Enabled)
        {
          DeleteOriginalFileIfNoPendingStatus(workspace, file, systemUser);
        }
			}
      catch (Exception e)
			{
        _logger.LogCritical(e, $"Failed on queue '{_coreSettings.PreviewWorker.DocumentStatusQueue}': {indentedMessage}");
        _documentStatusModel.BasicNack(args.DeliveryTag, false, false);
      }
		}
    
    private void OnTextExtractionStatus(object model, BasicDeliverEventArgs args)
    {
      string rawMessage = Encoding.UTF8.GetString(args.Body);
      PreviewStatusMessage message = JsonConvert.DeserializeObject<PreviewStatusMessage>(rawMessage);
      string indentedMessage = JsonConvert.SerializeObject(message, Formatting.Indented);
      _logger.LogTrace($"Received on queue '{_coreSettings.OcrWorker.TextExtractionStatusQueue}': {indentedMessage}");

      try
      {
        User systemUser = _userService.FindSystemUserAsync().Result;
        Workspace workspace = _workspaceService.FindAsync(message.WorkspaceId, systemUser).Result;
        InodeFacet inode = _inodeEngine.FindByIdAsync(message.InodeId, systemUser).Result;

        if (message.Status == "ready")
        {
          string text = System.IO.File.ReadAllText(_textExtractionService.GetPathAsync(workspace, inode, systemUser).Result);
          _inodeEngine.SetTextAsync(inode,text, systemUser).Wait();
        }

        var file = _fileService.FindLatestForInodeOrNullAsync(inode, systemUser).Result;
        _fileService.SetPropertyAsync(file, "file.ocr.textExtraction.status", message.Status, systemUser).Wait();
        
        _inodeNotificationService.SendInodesPropertiesUpdatedAsync(message.WorkspaceId, new[] {message.InodeId}).Wait();

        _logger.LogTrace($"Succeeded on queue '{_coreSettings.OcrWorker.TextExtractionStatusQueue}': {indentedMessage}");

        _textExtractionStatusModel.BasicAck(args.DeliveryTag, false);
        
        if (_coreSettings.S3Enabled)
        {
          DeleteOriginalFileIfNoPendingStatus(workspace, file, systemUser);
        }
      }
      catch (Exception e)
      {
        _logger.LogCritical(e, $"Failed on queue '{_coreSettings.OcrWorker.TextExtractionStatusQueue}': {indentedMessage}");
        _textExtractionStatusModel.BasicNack(args.DeliveryTag, false, false);
      }
    }
    
    private void OnSearchablePdfStatus(object model, BasicDeliverEventArgs args)
    {
      string rawMessage = Encoding.UTF8.GetString(args.Body);
      PreviewStatusMessage message = JsonConvert.DeserializeObject<PreviewStatusMessage>(rawMessage);
      string indentedMessage = JsonConvert.SerializeObject(message, Formatting.Indented);
      _logger.LogTrace($"Received on queue '{_coreSettings.OcrWorker.SearchablePdfStatusQueue}': {indentedMessage}");

      try
      {
        User systemUser = _userService.FindSystemUserAsync().Result;
        Workspace workspace = _workspaceService.FindAsync(message.WorkspaceId, systemUser).Result;
        Inode inode = _inodeEngine.FindByIdAsync(message.InodeId, systemUser).Result;
        
        var file = _fileService.FindLatestForInodeOrNullAsync(inode, systemUser).Result;

        if (message.Status == "ready")
        {
          _ocrQueueService.SendTextExtractionMessageAsync(workspace, inode, systemUser).Wait();
          _fileService.SetPropertyAsync(file, "file.ocr.textExtraction.status", "pending", systemUser).Wait();
        }
          
        _fileService.SetPropertyAsync(file, "file.ocr.searchablePdf.status", message.Status, systemUser).Wait();
        
        _inodeNotificationService.SendInodesPropertiesUpdatedAsync(message.WorkspaceId, new[] {message.InodeId}).Wait();

        _logger.LogTrace($"Succeeded on queue '{_coreSettings.OcrWorker.SearchablePdfStatusQueue}': {indentedMessage}");

        _searchablePdfStatusModel.BasicAck(args.DeliveryTag, false);
        
        if (_coreSettings.S3Enabled)
        {
          DeleteOriginalFileIfNoPendingStatus(workspace, file, systemUser);
        }
      }
      catch (Exception e)
      {
        _logger.LogCritical(e, $"Failed on queue '{_coreSettings.OcrWorker.SearchablePdfStatusQueue}': {indentedMessage}");
        _searchablePdfStatusModel.BasicNack(args.DeliveryTag, false, false);
      }
    }

    private void DeleteOriginalFileIfNoPendingStatus(Workspace workspace, File file, User systemUser)
    {
      string[] properties = {
        "file.imagePreview.status",
        "file.documentPreview.status",
        "file.tileMap.status",
        "file.ocr.textExtraction.status",
        "file.ocr.searchablePdf.status"
      };
      
      foreach (string property in properties)
      {
        try
        {
          if (_fileService.GetPropertyAsync(file, property, systemUser).Result == "pending")
          {
            return;
          }
        }
        catch
        {
          // ignored
        }
      }

      var message = new CleanupQueueMessage
      {
        Id = Guid.NewGuid().ToString(),
        Paths = new[] {_pathService.GetFileDirectory(workspace, file)}
      };
      SendQueueMessage(message, _coreSettings.CleanupQueue);
    }
  }
}