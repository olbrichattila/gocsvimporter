// Package main is running the CSV importer supporting large CSV files
package main

import (
	"flag"
	"os"

	importer "github.com/olbrichattila/gocsvimporter/internal"
)

func main() {
	help := flag.Bool("help", false, "Display help")

	// Parse the flags
	flag.Parse()

	// If help flag is set, display usage and exit
	if *help {
		displayHelp()
		os.Exit(0)
	}
	importer.Import()
}
