package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/starkandwayne/kafka-service-broker/cmd"
)

// Version set by ci/scripts/shipit; if "" it means "development"
var Version = ""

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" || os.Args[1] == "--version" {
			if Version == "" {
				fmt.Printf("kafka-service-broker (development)\n")
			} else {
				fmt.Printf("kafka-service-broker v%s\n", Version)
			}
			os.Exit(0)
		}
	}

	parser := flags.NewParser(&cmd.Opts, flags.Default)

	if len(os.Args) == 1 {
		_, err := parser.ParseArgs([]string{"--help"})
		if err != nil {
			os.Exit(1)
		}
	} else {
		_, err := parser.Parse()
		if err != nil {
			os.Exit(1)
		}
	}
}

func runBrk() {
}
