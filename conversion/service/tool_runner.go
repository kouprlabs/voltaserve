package service

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"

	"go.uber.org/zap"
)

type ToolRunner struct {
	logger *zap.SugaredLogger
}

func NewToolRunner() *ToolRunner {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &ToolRunner{
		logger: logger,
	}
}

func (r *ToolRunner) Run(inputPath string, opts core.ToolRunOptions) (outputPath string, stdout string, err error) {
	if inputPath != "" {
		for index, arg := range opts.Args {
			re := regexp.MustCompile(`\${input}`)
			substring := re.FindString(arg)
			if substring != "" {
				opts.Args[index] = re.ReplaceAllString(arg, inputPath)
			}
		}
	}
	for index, arg := range opts.Args {
		re := regexp.MustCompile(`\${output(?:\.[a-zA-Z0-9*#]+)*(?:\.[a-zA-Z0-9*#]+)?}`)
		substring := re.FindString(arg)
		if substring != "" {
			substring = regexp.MustCompile(`\${(.*?)}`).ReplaceAllString(substring, "$1")
			parts := strings.Split(substring, ".")
			if len(parts) == 1 {
				outputPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
				opts.Args[index] = re.ReplaceAllString(arg, outputPath)
			} else if len(parts) == 2 {
				outputPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + "." + parts[1])
				opts.Args[index] = re.ReplaceAllString(arg, outputPath)
			} else if len(parts) == 3 {
				if parts[1] == "*" {
					filename := filepath.Base(inputPath)
					outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
					if err := os.MkdirAll(outputDir, 0755); err != nil {
						return "", "", err
					}
					outputPath = filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + parts[2])
					opts.Args[index] = re.ReplaceAllString(arg, outputDir)
				} else if parts[1] == "#" {
					filename := filepath.Base(inputPath)
					basePath := filepath.FromSlash(os.TempDir() + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)))
					outputPath = filepath.FromSlash(basePath + "." + parts[2])
					opts.Args[index] = re.ReplaceAllString(arg, basePath)
				}
			}
		}
	}
	cmd := infra.NewCommand()
	r.logger.Named(infra.StrToolRunner).Infow("🔨  working", "bin", opts.Bin, "args", opts.Args)
	start := time.Now()
	if opts.Stdout {
		stdout, err := cmd.ReadOutput(opts.Bin, opts.Args...)
		elapsed := time.Since(start)
		if err != nil {
			r.logger.Named(infra.StrToolRunner).Errorw("⛈️  failed", "bin", opts.Bin, "args", opts.Args, "elapsed", elapsed, "error", "stdout", stdout, err.Error())
			return "", stdout, err
		} else {
			r.logger.Named(infra.StrToolRunner).Infow("🎉  succeeded", "bin", opts.Bin, "args", opts.Args, "elapsed", elapsed, "stdout", stdout)
			return outputPath, stdout, nil
		}
	} else {
		err := cmd.Exec(opts.Bin, opts.Args...)
		elapsed := time.Since(start)
		if err != nil {
			r.logger.Named(infra.StrToolRunner).Errorw("⛈️  failed", "bin", opts.Bin, "args", opts.Args, "elapsed", elapsed, "error", err.Error())
			return "", "", err
		} else {
			r.logger.Named(infra.StrToolRunner).Infow("🎉  succeeded", "bin", opts.Bin, "args", opts.Args, "elapsed", elapsed)
			return outputPath, "", err
		}
	}
}
