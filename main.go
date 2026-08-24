// Command buoy-qc provides buoy data QC and sea-state analysis as both a CLI
// tool and an HTTP service.
package main

import (
	"fmt"
	"os"

	"buoy-qc/internal/cli"
	"buoy-qc/internal/server"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "serve" {
		runServer(args[1:])
		return
	}
	if len(args) == 0 {
		runServer(nil)
		return
	}
	os.Exit(cli.Run(args, os.Stdout, os.Stderr))
}

func runServer(args []string) {
	addr := ":8080"
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-addr" || args[i] == "--addr" {
			addr = args[i+1]
			break
		}
	}
	cfg := server.Config{Addr: addr}
	fmt.Fprintf(os.Stdout, "buoy-qc server listening on %s\n", server.FormatAddr(addr))
	if err := server.ListenAndServe(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
