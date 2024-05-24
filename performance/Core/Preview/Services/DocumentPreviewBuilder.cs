namespace Defyle.Core.Preview.Services
{
  using System;
  using System.Diagnostics;
  using System.IO;
  using System.Threading;
  using Infrastructure.Poco;
  using Queue;

  public class DocumentPreviewBuilder
  {
    private readonly PreviewMessage.DocumentPreviewPayload _options;

    private readonly ExternalAppSettings _binary;

    public DocumentPreviewBuilder(
      PreviewMessage.DocumentPreviewPayload options,
      ExternalAppSettings binary)
    {
      _options = options;
      _binary = binary;
    }

    public void Build()
    {
      bool requiresCleanup = false;
      
      if (!Directory.Exists(_options.OutputDirectory))
      {
        Directory.CreateDirectory(_options.OutputDirectory);
        requiresCleanup = true;
      }

      try
      {
        if (!File.Exists(_options.File))
        {
          throw new Exception($"Input file '{_options.File}' does not exist.");
        }

        string binary;
        string args;

        if (_binary.Docker != null)
        {
          var docker = _binary.Docker;

          binary = docker.Cli;
          args = $"run --rm -it " +
                 $"-v {Path.GetDirectoryName(_options.File)}:{docker.InputDirectory} " +
                 $"-v {_options.OutputDirectory}:{docker.OutputDirectory} " +
                 $"{docker.Image} " +
                 $"{_binary.Docker.App} " +
                 "--headless " +
                 "--convert-to pdf " +
                 $"{Path.Combine(docker.InputDirectory, Path.GetFileName(_options.File))} " +
                 $"--outdir {docker.OutputDirectory}";
        }
        else
        {
          binary = _binary.App;
          args = $"--headless --convert-to pdf {_options.File} --outdir {_options.OutputDirectory}";
        }

        string command = $"{binary} {args}";

        var processInfo = new ProcessStartInfo(binary, args);
        processInfo.CreateNoWindow = true;
        processInfo.UseShellExecute = false;

        using (var process = new Process())
        {
          process.StartInfo = processInfo;

          process.Start();
          process.WaitForExit(300000);

          if (!process.HasExited)
          {
            process.Kill();
          }

          if (process.ExitCode != 0)
          {
            throw new Exception($"Command {command} failed with status code {process.ExitCode}.");
          }
        }

        string outputFile = Path.Combine(_options.OutputDirectory, Path.GetFileNameWithoutExtension(_options.File) + ".pdf");

        if (!File.Exists(outputFile))
        {
          Thread.Sleep(5000);
        }

        if (!File.Exists(outputFile))
        {
          throw new Exception($"Output file '{outputFile}' does not exist.");
        }

        File.Move(outputFile, Path.Combine(_options.OutputDirectory, "document.pdf"));
      }
      catch (Exception)
      {
        if (requiresCleanup)
        {
          Directory.Delete(_options.OutputDirectory, true);
        }
        throw;
      }
    }
  }
}