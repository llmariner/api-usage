package main

import (
	"log"
	"os"
)

func main() {
	cmd := runCmd()
	cmd.SilenceUsage = true

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
