// Command buoy-qc is a CLI for buoy data QC and sea-state analysis.
package main

import (
	"os"

	"buoy-qc/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
