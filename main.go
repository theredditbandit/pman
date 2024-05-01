package main

import (
	"log"

	"github.com/theredditbandit/pman/cmd"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	cmd.Execute()
}
