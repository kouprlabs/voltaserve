namespace Defyle.Core.Infrastructure.Services
{
  using System;
  using System.Text;
  using Newtonsoft.Json;
  using Poco;
  using RabbitMQ.Client;

  public class QueueService
  {
    private readonly MessageBrokerSettings _settings;

    protected QueueService(MessageBrokerSettings settings)
    {
      _settings = settings;
    }
    
    public void SendQueueMessage(object message, string queue)
    {
      var factory = new ConnectionFactory
      {
        Uri = new Uri(_settings.Url)
      };

      using var connection = factory.CreateConnection();
      using var channel = connection.CreateModel();
      channel.QueueDeclare(queue, true, false, false, null);

      var properties = channel.CreateBasicProperties();
      properties.Persistent = true;
          
      channel.BasicPublish(string.Empty, queue, properties,
        Encoding.UTF8.GetBytes(JsonConvert.SerializeObject(message))
      );
    }
  }
}