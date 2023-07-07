package service

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
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
			return "", stdout, err
		} else {
			return outputFile, stdout, nil
		}
	} else {
		if err := cmd.Exec(opts.Bin, opts.Args...); err != nil {
			return "", "", err
		} else {
			return outputFile, "", err
		}
	}
}
