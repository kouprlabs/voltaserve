namespace Defyle.Core.Infrastructure.Poco
{
  using System;
  using System.IO;

  public static class PathUtils
  {
    public static string Rewrite(string value)
    {
      if (!string.IsNullOrWhiteSpace(value) && !Path.IsPathFullyQualified(value) && value.StartsWith("~/"))
      {
        return Path.Combine(
          Environment.GetFolderPath(Environment.SpecialFolder.UserProfile),
          value.Replace("~/", string.Empty));
      }
      else
      {
        return value;
      }
    }
  }
}