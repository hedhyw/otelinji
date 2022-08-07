package main

import (
	"log"
	"os"

	"github.com/hedhyw/otelinj/internal/app"
	"github.com/hedhyw/otelinj/internal/pkg/config"
)

func main() {
	cfg, err := config.FromCLI(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = app.New(cfg).Run(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
