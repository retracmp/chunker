package chunk

import (
	"acid/chunker/src/helpers"
	"strings"

	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Chunker struct {
	ID string
	RootPath string
	Files []*File
	ChunkSize int64
}

func NewChunker(rootPath string, chunkSize int64) *Chunker {
	return &Chunker{
		ID: strings.ToUpper(helpers.MD5([]byte(filepath.Base(rootPath)))),
		RootPath: rootPath,
		ChunkSize: chunkSize,
		Files: make([]*File, 0),
	}
}

func (c *Chunker) Chunk() error {
	paths, err := Paths(c.RootPath)
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err = c.processFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chunker) processFile(path string) error {
	file, err := NewFile(path, c.RootPath)
	if err != nil {
		return err
	}
	if err = file.Chunk(c.ChunkSize, c.ID); err != nil {
		return err
	}

	fmt.Printf("File::%s::Size::%d\n", file.Path, file.Size)

	c.Files = append(c.Files, file)

	return nil
}

func (c *Chunker) ChunkThreaded() error {
	paths, err := Paths(c.RootPath)
	if err != nil {
		return err
	}

	wait := sync.WaitGroup{}
	limiter := make(chan struct{}, 15)

	for _, path := range paths {
		limiter <- struct{}{}
		wait.Add(1)
		go c.processFileThreaded(path, &wait, limiter)
	}
	wait.Wait()

	return nil
}

func (c *Chunker) processFileThreaded(path string, wait *sync.WaitGroup, limiter chan struct{}) {
	defer wait.Done()
	defer func() { <-limiter }()

	file, err := NewFile(path, c.RootPath)
	if err != nil {
		return
	}
	if err = file.Chunk(c.ChunkSize, c.ID); err != nil {
		return
	}

	fmt.Printf("File::%s::Size::%d \n", file.Path, file.Size)

	c.Files = append(c.Files, file)
}

type RenderedChunks struct {
	ID string
	UploadName string
	Files []*File
}

func (c *Chunker) render() []byte {
	renderedChunks := &RenderedChunks{
		ID: c.ID,
		UploadName: filepath.Base(c.RootPath),
		Files: c.Files,
	}

	return helpers.JSONToBytes(renderedChunks)
}

func (c *Chunker) RenderToFile() error {
	bytes := c.render()

	file, err := os.Create(fmt.Sprintf("./builds/%s.acidmanifest", filepath.Base(c.RootPath)))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	fmt.Printf("Rendered to: %s \n", fmt.Sprintf("./builds/%s.acidmanifest", filepath.Base(c.RootPath)))

	return nil
}