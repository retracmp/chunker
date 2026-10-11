package chunk

import (
	"acid/chunker/src/chunk"
	"acid/chunker/src/helpers"
	"os"
	"strings"

	"fmt"
)

func ArgOrPanic(args []string, index int) string {
	if len(args) <= index {
		panic(fmt.Sprintf("Missing argument at index %d", index))
	}
	return args[index]
}

// chunk <root> [name] [manifest filename] [-WL:<file>]... [-DEL:<path>]...
// -WL: limits chunking to the named files; -DEL: adds a tombstone so a launcher applying this layer removes the path.
func Start() error {
	t := helpers.NewPerformanceTimer()
	defer t.EndTimer()

	positional := make([]string, 0)
	whitelist := make([]string, 0)
	deleted := make([]string, 0)
	
	for _, arg := range os.Args[2:] {
		switch {
		case strings.HasPrefix(arg, "-WL:"):
			whitelist = append(whitelist, strings.TrimPrefix(arg, "-WL:"))
		case strings.HasPrefix(arg, "-DEL:"):
			deleted = append(deleted, strings.TrimPrefix(arg, "-DEL:"))
		default:
			positional = append(positional, arg)
		}
	}

	chunker := chunk.NewChunker(ArgOrPanic(positional, 0), 128 * chunk.MB)
	if len(positional) > 1 {
		chunker.SetCustomName(positional[1])
	}
	for _, name := range whitelist {
		chunker.AddFileToWhitelist(name)
	}
	for _, path := range deleted {
		chunker.AddDeletedFile(path)
	}

	if err := chunker.Chunk(); err != nil {
		return fmt.Errorf("chunker.Chunk: %w", err)
	}
	if err := chunker.RenderToFile(positional[2:]...); err != nil {
		return fmt.Errorf("chunker.RenderToFile: %w", err)
	}

	return nil
}