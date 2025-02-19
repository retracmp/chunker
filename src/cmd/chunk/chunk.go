package chunk

import (
	"acid/chunker/src/chunk"
	"acid/chunker/src/helpers"
	"os"

	"fmt"
)

func ArgOrPanic(args []string, index int) string {
	if len(args) <= index {
		panic(fmt.Sprintf("Missing argument at index %d", index))
	}
	return args[index]
}

func Start() error {
	t := helpers.NewPerformanceTimer()
	defer t.EndTimer()

	chunker := chunk.NewChunker(ArgOrPanic(os.Args[2:], 0), 128 * chunk.MB)
	if err := chunker.Chunk(); err != nil {
		return fmt.Errorf("chunker.Chunk: %w", err)
	}
	if err := chunker.RenderToFile(); err != nil {
		return fmt.Errorf("chunker.RenderToFile: %w", err)
	}

	return nil
}