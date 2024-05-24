namespace Defyle.Core.Storage.Models
{
  using System;
  using Microsoft.AspNetCore.StaticFiles;

  public class File
	{
		private static readonly FileExtensionContentTypeProvider FileExtensionContentTypeProvider = new FileExtensionContentTypeProvider();

    public File()
    {
    }
    
		public File(string extension, long size, bool indexContent)
		{
			Id = Guid.NewGuid().ToString().Replace("-", "");
      Extension = extension;
			Size = size;
      IndexContent = indexContent;
      Mime = GetMime();
    }

    public string Id { get; set; }

    public string Extension { get; set; }
    
		public long Size { get; set; }

    public bool IndexContent { get; set; }

    public string Mime { get; set; }

    public long Version { get; set; }
    
    public long CreateTime { get; set; }
    
    public long UpdateTime { get; set; }

		public string GetMime()
		{
			if (Extension == "csv")
			{
				return "text/csv";
			}

			try
			{
        switch (Extension)
        {
          case "odt":
          case "fodt":
            return "application/vnd.oasis.opendocument.text";
          
          case "odp":
          case "fodp":
            return "application/vnd.oasis.opendocument.presentation";
          
          case "ods":
          case "fods":
            return "application/vnd.oasis.opendocument.spreadsheet";
          
          default:
            return FileExtensionContentTypeProvider.Mappings["." + Extension];
        }
      }
			catch (Exception)
			{
				return "application/octet-stream";
			}
		}
  }
}