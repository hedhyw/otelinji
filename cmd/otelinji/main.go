package main

import (
	"log"
	"os"

	"github.com/hedhyw/otelinji/internal/app"
	"github.com/hedhyw/otelinji/internal/pkg/config"
)

// Version will be set on build.
var version = "unknown"

func main() {
	cfg, err := config.FromCLI(os.Args[1:], version)
	if err != nil {
		log.Fatal(err)
	}

	err = app.New(cfg).Run(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
