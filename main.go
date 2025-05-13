package main

import (
	"acid/chunker/src/cmd"

	"os"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "ui")
	}

	if len(os.Args) < 2 {
		return
	}

	err := cmd.RunService(cmd.ChunkerCommand(os.Args[1]))
	if err != nil {
		panic(err)
	} 
}