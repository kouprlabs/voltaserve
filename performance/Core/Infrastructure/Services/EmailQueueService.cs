namespace Defyle.Core.Infrastructure.Services
{
  using Poco;

  public class EmailQueueService : QueueService
  {
    private readonly CoreSettings _coreSettings;

    public EmailQueueService(CoreSettings coreSettings)
      : base(coreSettings.MessageBroker)
    {
      _coreSettings = coreSettings;
    }

    public void SendEmailMessage(EmailMessage message)
    {
      SendQueueMessage(message, _coreSettings.EmailQueue);
    }
  }
}