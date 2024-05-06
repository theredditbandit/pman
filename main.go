package main

import (
	"log"
	"os"

	"github.com/theredditbandit/pman/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
