// Package main is running the CSV importer supporting large CSV files
package main

import (
	importer "github.com/olbrichattila/gocsvimporter/internal"
)

func main() {
	importer.Import()
}
