namespace Defyle.Core.Inode.Services
{
  using System;
  using System.Collections.Generic;
  using System.Text;
  using System.Threading.Tasks;
  using Auth.Models;
  using Infrastructure.Poco;
  using Microsoft.Extensions.Logging;
  using Models;
  using Newtonsoft.Json;
  using Pocos;
  using RabbitMQ.Client;
  using RabbitMQ.Client.Events;
  using Workspace.Models;
  using Workspace.Services;

  public class SftpService
	{
    private readonly CoreSettings _coreSettings;
    private readonly WorkspaceNotificationService _workspaceNotificationService;
    private readonly ILogger<SftpService> _logger;
    private ConnectionFactory _factory;
		private IConnection _connection;
		private IModel _channel;
		private string _replyQueueName;
		private EventingBasicConsumer _consumer;

		public SftpService(
			CoreSettings coreSettings,
      WorkspaceNotificationService workspaceNotificationService,
      ILogger<SftpService> logger)
    {
			_coreSettings = coreSettings;
      _workspaceNotificationService = workspaceNotificationService;
      _logger = logger;
    }

		public void Register()
		{
			_factory = new ConnectionFactory
			{
        Uri = new Uri(_coreSettings.MessageBroker.Url)
			};

			_connection = _factory.CreateConnection();
			_channel = _connection.CreateModel();
			_channel.QueueDeclare(_coreSettings.SftpQueue, true, false, false, null);

			_replyQueueName = _channel.QueueDeclare().QueueName;

			_consumer = new EventingBasicConsumer(_channel);
			_consumer.Received += (model, ea) =>
			{
				string rawResponse = Encoding.UTF8.GetString(ea.Body);
				SftpImportResponseMessage response = JsonConvert.DeserializeObject<SftpImportResponseMessage>(rawResponse);
				string indentedMessage = JsonConvert.SerializeObject(response, Formatting.Indented);
				_logger.LogTrace($"Response: {indentedMessage}");

				// TODO: remove task where Id == response.MessageId
        throw new NotImplementedException();

        _workspaceNotificationService.SendWorkspaceTasksUpdatedAsync(response.WorkspaceId).Wait();
        
				string title;
				List<string> content;
				if (response.Success)
				{
          title = "SFTP Import";
          content = new List<string> {"Succeeded 🌈"};
        }
				else
				{
          title = "SFTP Import";
          content = new List<string> {"Failed ⛈", response.StatusMessage};
        }
        
        _workspaceNotificationService.SendWorkspaceNotificationAsync(response.WorkspaceId, title, content).Wait();
      };

			_channel.BasicConsume(
				consumer: _consumer,
				queue: _replyQueueName,
				autoAck: true);
		}

		public void Deregister()
		{
			_channel.Close();
		}

		public async Task ImportAsync(User user, Workspace workspace, Inode inode, SftpImportOptions request, string cipherPassword)
		{
      SftpImportMessage message = new SftpImportMessage();
			message.Id = Guid.NewGuid().ToString();
			message.WorkspaceId = workspace.Id;
			message.IsWorkspaceEncrypted = workspace.Encrypted;
			message.ParentNodeId = inode.Id;
			message.UserId = user.Id;
			message.Request = request;
			message.TempDirectory = _coreSettings.TempDirectory;
			message.CipherPassword = cipherPassword;
			message.InheritAccessRights = true;
      message.ApiUrl = _coreSettings.Url;
			message.PartitionId = _coreSettings.PartitionId;

			WorkspaceTask task = new WorkspaceTask
			{
				Id = message.Id,
				Type = WorkspaceTask.WorkspaceTaskTypeIndeterminateProgress,
				Title = "SFTP Import",
				Content = "In progress...",
        CreatedAt = DateTimeOffset.UtcNow.ToUnixTimeMilliseconds()
			};
      
      // TODO: append task
      throw new NotImplementedException();
      
      await _workspaceNotificationService.SendWorkspaceTasksUpdatedAsync(workspace.Id);

			var properties = _channel.CreateBasicProperties();
			properties.Persistent = true;
			properties.CorrelationId = Guid.NewGuid().ToString();
			properties.ReplyTo = _replyQueueName;

			_channel.BasicPublish(
				string.Empty,
				_coreSettings.SftpQueue,
				properties,
				Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(message))
			);
		}
	}
}