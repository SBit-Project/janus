package main

import (
	"log"

	"github.com/SBit-Project/janus/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cli.Run()
}
