package chunk

import (
	"acid/chunker/src/helpers"

	"fmt"
	"os"
	"strings"
)

// compressing takes 10x longer to process, on both uploading and compiling the parts
// however, it does reduce the size of the files by ~50% is most cases so it's worth it
const compress = true

type Chunk struct {
	Offset int64
	Size int64
	Hash string
}

func NewChunk(file *File, offset int64, size int64) *Chunk {
	if offset+size > file.Size {
		size = file.Size - offset
	}

	return &Chunk{
		Offset: offset,
		Size: size,
	}
}

func (c *Chunk) Process(file *os.File, buildId string) error {
	bytes := make([]byte, c.Size)
	
	if _, err := file.Seek(c.Offset, 0); err != nil {
		return err
	}
	if _, err := file.Read(bytes); err != nil {
		return err
	}

	if compress {
		bytes_, err := helpers.Compress(bytes)
		if err != nil {
			return err
		}
		bytes = bytes_
	}

	c.Hash = strings.ToUpper(helpers.MD5(bytes))

	file, err := os.Create(fmt.Sprintf("./builds/%s/%s", buildId, c.Hash))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fmt.Sprintf("./builds/%s", buildId), 0755); err != nil {
			return err	
		}

		file, err = os.Create(fmt.Sprintf("./builds/%s/%s", buildId, c.Hash))
	}
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (c *Chunk) Data(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("os.ReadFile: %w", err)
		return nil
	}

	if compress {
		file_, err := helpers.Decompress(file)
		if err != nil {
			fmt.Println("helpers.Decompress: %w", err)
			return nil
		}
		file = file_
	}

	return file
}