namespace Defyle.Core.Ocr.Services
{
  using System;
  using System.Diagnostics;
  using System.IO;
  using System.Threading;
  using Exceptions;
  using Infrastructure.Poco;
  using Queue;

  public class TextExtractionBuilder
  {
    private readonly TextExtractionMessage.TextExtractionPayload _options;

    private readonly ExternalAppSettings _binary;
    private readonly long _maxTextLength;
    private readonly long _maxFileSize;

    public TextExtractionBuilder(
      TextExtractionMessage.TextExtractionPayload options,
      ExternalAppSettings binary,
      long maxTextLength,
      long maxFileSize)
    {
      _options = options;
      _binary = binary;
      _maxTextLength = maxTextLength;
      _maxFileSize = maxFileSize;
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

        const string filename = "document.txt";

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
                 $"{Path.Combine(docker.InputDirectory, Path.GetFileName(_options.File))} " +
                 $"{Path.Combine(docker.OutputDirectory, filename)}";
        }
        else
        {
          binary = _binary.App;
          args = $"{_options.File} {Path.Combine(_options.OutputDirectory, filename)}";
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

        string outputFile = Path.Combine(_options.OutputDirectory, filename);

        if (!File.Exists(outputFile))
        {
          Thread.Sleep(5000);
        }

        if (!File.Exists(outputFile))
        {
          throw new Exception($"Output file '{outputFile}' does not exist.");
        }

        var fileInfo = new FileInfo(outputFile);
        if (fileInfo.Length > _maxFileSize)
        {
          throw new TextFileTooLargeException();
        }

        string text = File.ReadAllText(outputFile);
        if (text.Length > _maxTextLength)
        {
          throw new TextFileTooLargeException();
        }
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