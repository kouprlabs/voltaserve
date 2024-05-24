namespace Defyle.WebApi.Configuration
{
  using System;
  using Core.Infrastructure.Poco;
  using Core.Inode.Models;
  using Microsoft.Extensions.DependencyInjection;
  using Microsoft.Extensions.Logging;
  using Nest;

  public static class ElasticsearchExtensions
  {
    public static void AddElasticsearch(this IServiceCollection services, ElasticsearchSettings elasticsearchSettings,
      ILogger<Startup> logger)
    {
      var settings = new ConnectionSettings(new Uri(elasticsearchSettings.Url))
        .DefaultMappingFor<Inode>(m => m
          .IndexName(elasticsearchSettings.InodesIndex)
          .IdProperty(p => p.Id)
        );

      var client = new ElasticClient(settings);
      logger.LogInformation("Elasticsearch connected.");

      if (!client.Indices.Exists(elasticsearchSettings.InodesIndex).Exists)
      {
        var response = client.Indices.Create(elasticsearchSettings.InodesIndex, index => index.Map<Inode>(x => x.AutoMap()));
        if (!response.IsValid)
        {
          logger.LogCritical(response.ToString());
        }
      }

      services.AddSingleton<IElasticClient>(client);
    }
  }
}