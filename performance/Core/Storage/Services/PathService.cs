namespace Defyle.Core.Storage.Services
{
  using System.IO;
  using Infrastructure.Poco;
  using Workspace.Models;
  using File = Models.File;

  public class PathService
	{
		private readonly CoreSettings _coreSettings;

		public PathService(CoreSettings coreSettings)
		{
			_coreSettings = coreSettings;
		}

		public string GetOriginalFile(Workspace workspace, File file)
		{
			string path = Path.Combine(
				GetFileDirectory(workspace, file),
				$"original.{file.Extension}");

			if (workspace.Encrypted)
			{
				path += ".enc";
			}

			return path;
		}

    public string GetFileDirectory(Workspace workspace, File file) =>
      Path.Combine(_coreSettings.DataDirectory, workspace.Id, "file-data", file.Id);

    public string GetS3OriginalFile(Workspace workspace, File file)
    {
      string path = Path.Combine(GetS3FileDataDirectory(workspace, file), $"original.{file.Extension}");
      if (workspace.Encrypted)
      {
        path += ".enc";
      }

      return path;
    }

    private string GetS3FileDataDirectory(Workspace workspace, File file) =>
      Path.Combine(workspace.Id, "file-data", file.Id);
  }
}