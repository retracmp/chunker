package cmd

import (
	"acid/chunker/src/cmd/chunk"
	"acid/chunker/src/cmd/compile"
	"acid/chunker/src/cmd/download"
	"acid/chunker/src/cmd/hoster"

	"fmt"
)

type ChunkerCommand string
const (
	Chunk ChunkerCommand = "chunk"
	Compile ChunkerCommand = "compile"
	Download ChunkerCommand = "download"
	Hoster ChunkerCommand = "hoster"
)

var (
	commandLookup = map[ChunkerCommand]func()error{
		Chunk: chunk.Start,
		Compile: compile.Start,
		Download: download.Start,
		Hoster: hoster.Start,
	}
)

func RunService(t ChunkerCommand) error {
	if handler := commandLookup[t]; handler != nil {
		if err := handler(); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("no service found with id '%s'", t)
}