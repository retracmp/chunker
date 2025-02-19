package compile

import (
	"acid/chunker/src/chunk"
	"acid/chunker/src/helpers"
	"os"
	"path"
	"strings"
)

type Compiler struct {
	ManifestPath string
	ChunksPath string
	CompilePath string

	Files []*chunk.File
}

func NewCompiler(manifestPath string, chunksPath string, compilePath string) *Compiler {
	manifest := helpers.JSONFromFile[chunk.RenderedChunks](manifestPath)

	return &Compiler{
		ManifestPath: manifestPath,
		ChunksPath: path.Join(chunksPath, manifest.ID),
		CompilePath: compilePath,
		Files: manifest.Files,
	}
}

func (c *Compiler) Compile() error {
	for _, file := range c.Files {
		file.DisplayPath = strings.ReplaceAll(file.DisplayPath, "\\\\", "\\")
		file.DisplayPath = strings.ReplaceAll(file.DisplayPath, "\\", "/")
		if err := c.processFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) Check() error {
	for _, file := range c.Files {
		file.DisplayPath = strings.ReplaceAll(file.DisplayPath, "\\\\", "\\")
		file.DisplayPath = strings.ReplaceAll(file.DisplayPath, "\\", "/")
		if err := c.checkFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) processFile(file *chunk.File) error {
	if err := file.Rebuild(c.ChunksPath, c.CompilePath); err != nil {
		return err
	}

	return nil
}

func (c *Compiler) checkFile(file *chunk.File) error {
	if err := file.Check(c.ChunksPath, c.CompilePath); err != nil {
		return err
	}

	return nil
}

func (c *Compiler) Cleanup(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	return nil
}