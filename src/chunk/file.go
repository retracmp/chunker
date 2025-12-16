package chunk

import (
	"acid/chunker/src/helpers"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type File struct {
	Name string
	Path string `json:"-"`
	DisplayPath string `json:"Path"`
	Hash string
	Size int64
	Chunks []*Chunk
}

func NewFile(path string, rootPath string) (*File, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat: %w", err)
	}

	return &File{
		Name: info.Name(),
		Path: path,
		Size: info.Size(),
		DisplayPath: strings.TrimPrefix(path, rootPath),
	}, nil
}

func (f *File) Chunk(chunkSize int64, buildId string) error {
	file, err := os.OpenFile(f.Path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes := make([]byte, f.Size)
	if _, err := io.ReadFull(file, bytes); err != nil {
		return err
	}
	f.Hash = strings.ToUpper(helpers.MD5(bytes))
	bytes = nil // free 

	offset := int64(0)
	for offset < f.Size {
		chunk := NewChunk(f, offset, chunkSize)
		if err := chunk.Process(file, buildId); err != nil {
			return err
		}
		
		f.Chunks = append(f.Chunks, chunk)

		fmt.Printf("Chunk::%s::Size::%d \n", chunk.Hash, chunk.Size)

		offset += chunk.Size
	}

	return nil
}

func (f *File) Rebuild(chunkPath string, compilePath string) error {
	if err := os.MkdirAll(path.Dir(path.Join(compilePath, f.DisplayPath)), os.ModePerm); err != nil {
		return err
	}

	if f.alreadyBuilt(compilePath) {
		fmt.Printf("RebuiltFile::%s::Size::%d \n", f.DisplayPath, f.Size)
		return nil
	}

	file, err := os.Create(path.Join(compilePath, f.DisplayPath))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, chunk := range f.Chunks {
		if _, err := file.Write(chunk.Data(fmt.Sprintf("%s/%s", chunkPath, chunk.Hash))); err != nil {
			return err
		}
	}

	fmt.Printf("RebuiltFile::%s::Size::%d \n", f.DisplayPath, f.Size)

	return nil
}

func (f *File) Check(chunkPath string, compilePath string) error {
	if err := os.MkdirAll(path.Dir(path.Join(compilePath, f.DisplayPath)), os.ModePerm); err != nil {
		return err
	}

	if f.alreadyBuilt(compilePath) {
		fmt.Printf("CheckFile::%s::Size::%d \n", f.DisplayPath, f.Size)
		return nil
	}

	fmt.Printf("CheckFile::%s::Size::%d \n", f.DisplayPath, f.Size)

	return fmt.Errorf("File does not exist")
}

func (f *File) alreadyBuilt(compilePath string) bool {
	info, err := os.Stat(path.Join(compilePath, f.DisplayPath))
	if os.IsNotExist(err) {
		return false
	}

	return info.Size() == f.Size
}