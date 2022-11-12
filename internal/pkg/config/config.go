package config

import (
	"flag"
	"fmt"
)

// Config of the application.
type Config struct {
	FileName         string
	TemplateName     string
	WriteIntoFile    bool
	SkipGenerated    bool
	Version          string
	OnlyPrintVersion bool
}

// FromCLI parses command line arguments and returns a config.
func FromCLI(args []string, version string) (*Config, error) {
	flagSet := flag.NewFlagSet("otelinji", flag.ContinueOnError)

	fileName := flagSet.String(
		"filename",
		"",
		"golang file [required]",
	)

	onlyPrintVersion := flagSet.Bool(
		"version",
		false,
		"print application version and quit",
	)

	templateName := flagSet.String(
		"template",
		"@/otel",
		"path to template file [optional]",
	)

	writeIntoFile := flagSet.Bool(
		"w",
		false,
		"write result to file [optional]",
	)

	skipGenerated := flagSet.Bool(
		"skip-generated",
		true,
		"skip files with `DO NOT EDIT` comment",
	)

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing command-line: %w", err)
	}

	return &Config{
		FileName:         *fileName,
		TemplateName:     *templateName,
		WriteIntoFile:    *writeIntoFile,
		SkipGenerated:    *skipGenerated,
		Version:          version,
		OnlyPrintVersion: *onlyPrintVersion,
	}, nil
}
