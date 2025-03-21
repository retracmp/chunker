package main

import (
	"acid/chunker/src/cmd"
	"fmt"

	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command>\n", os.Args[0])
		return
	}

	err := cmd.RunService(cmd.ChunkerCommand(os.Args[1]))
	if err != nil {
		panic(err)
	} 
}