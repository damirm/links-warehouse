package main

import (
	"flag"
	"log"
	"os"

	"github.com/damirm/links-warehouse/internal/command"
)

func main() {
	cmd := command.NewWarehouseCommand()

	fs := flag.NewFlagSet("warehouse", flag.ExitOnError)
	cmd.ExportFlags(fs)

	args := os.Args[1:]

	if err := fs.Parse(args); err != nil {
		log.Printf("failed to parse args: %v", err)
		os.Exit(1)
	}

	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
