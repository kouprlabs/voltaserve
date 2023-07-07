package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var StrRunner = fmt.Sprintf("%-13s", "runner")

type Runner struct {
	logger *zap.SugaredLogger
}

func NewRunner() *Runner {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableCaller = true
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return &Runner{
		logger: logger.Sugar(),
	}
}

func (r *Runner) Run(inputPath string, opts core.RunOptions) (outputFile string, stdout string, err error) {
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
				outputFile = filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
				opts.Args[index] = re.ReplaceAllString(arg, outputFile)
			} else if len(parts) == 2 {
				outputFile = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + "." + parts[1])
				opts.Args[index] = re.ReplaceAllString(arg, outputFile)
			} else if len(parts) == 3 {
				if parts[1] == "*" {
					filename := filepath.Base(inputPath)
					outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
					if err := os.MkdirAll(outputDir, 0755); err != nil {
						return "", "", err
					}
					outputFile = filepath.FromSlash(outputDir + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + parts[2])
					opts.Args[index] = re.ReplaceAllString(arg, outputDir)
				} else if parts[1] == "#" {
					filename := filepath.Base(inputPath)
					basePath := filepath.FromSlash(os.TempDir() + "/" + strings.TrimSuffix(filename, filepath.Ext(filename)))
					outputFile = filepath.FromSlash(basePath + "." + parts[2])
					opts.Args[index] = re.ReplaceAllString(arg, basePath)
				}
			}
		}
	}
	cmd := infra.NewCommand()
	if opts.Stdout {
		stdout, err := cmd.ReadOutput(opts.Bin, opts.Args...)
		if err != nil {
			r.logger.Named(StrRunner).Errorw("⛈️  failed", "bin", opts.Bin, "args", opts.Args, "error", "stdout", stdout, err.Error())
			return "", stdout, err
		} else {
			r.logger.Named(StrRunner).Infow("🎉  succeeded", "bin", opts.Bin, "args", opts.Args, "stdout", stdout)
			return outputFile, stdout, nil
		}
	} else {
		if err := cmd.Exec(opts.Bin, opts.Args...); err != nil {
			r.logger.Named(StrRunner).Errorw("⛈️  failed", "bin", opts.Bin, "args", opts.Args, "error", err.Error())
			return "", "", err
		} else {
			r.logger.Named(StrRunner).Infow("🎉  succeeded", "bin", opts.Bin, "args", opts.Args)
			return outputFile, "", err
		}
	}
}
