namespace Defyle.Core.Infrastructure.Services
{
  using Newtonsoft.Json;
  using Newtonsoft.Json.Serialization;

  public class JsonService
  {
    private JsonSerializerSettings _serializerSettings;
    
    public JsonSerializerSettings GetJsonSerializerSettings()
    {
      if (_serializerSettings == null)
      {
        _serializerSettings = new JsonSerializerSettings();
        _serializerSettings.ContractResolver = new CamelCasePropertyNamesContractResolver();
      }

      return _serializerSettings;
    }
  }
}