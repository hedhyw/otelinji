package config

import (
	"flag"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

// Config of the application.
type Config struct {
	FileName      string
	TemplateName  string
	WriteIntoFile bool
	SkipGenerated bool
}

const (
	errFileNameRequired semerr.Error = "required `-filename` is not specified"
)

// FromCLI parses command line arguments and returns a config.
func FromCLI(args []string) (*Config, error) {
	fileName := flag.CommandLine.String(
		"filename",
		"",
		"golang file [required]",
	)

	templateName := flag.CommandLine.String(
		"template",
		"@/otel",
		"path to template file [optional]",
	)

	writeIntoFile := flag.CommandLine.Bool(
		"w",
		false,
		"write result to file [optional]",
	)

	skipGenerated := flag.CommandLine.Bool(
		"skip-generated",
		true,
		"skip files with `DO NOT EDIT` comment",
	)

	if err := flag.CommandLine.Parse(args); err != nil {
		return nil, fmt.Errorf("parsing command-line: %w", err)
	}

	if *fileName == "" {
		return nil, errFileNameRequired
	}

	return &Config{
		FileName:      *fileName,
		TemplateName:  *templateName,
		WriteIntoFile: *writeIntoFile,
		SkipGenerated: *skipGenerated,
	}, nil
}
