package downloader

import (
	"acid/chunker/src/chunk"
	"acid/chunker/src/helpers"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Downloader struct {
	Manifest *chunk.RenderedChunks
	client *http.Client
	root string
	path string
	compilePath string
	
	maxBytesPerSecond int64
}

func NewDownloader(url string, downloadPath string, compilePath string) *Downloader {
	if url[len(url)-1:] == "/" {
		url = url[:len(url)-1]
	}

	return &Downloader{
		client: &http.Client{
			Timeout: 0,
		},
		root: url,
		path: downloadPath,
		compilePath: compilePath,
	}
}

func (d *Downloader) FetchManifest(manifest string, fileOutput string) (*chunk.RenderedChunks, error) {
	resp, err := d.get(fmt.Sprintf("%s/%s", d.root, manifest))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := helpers.JSONFromBytes[chunk.RenderedChunks](bytes)
	if result.ID == "" {
		return nil, fmt.Errorf("failed to parse manifest")
	}
	d.Manifest = &result

	file, err := os.Create(fmt.Sprintf("%s/%s", fileOutput, manifest))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fileOutput, os.ModePerm); err != nil {
			return nil, err
		}

		file, err = os.Create(fmt.Sprintf("%s/%s", fileOutput, manifest))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *Downloader) Download() error {
	for _, file := range d.Manifest.Files {
		fmt.Printf("File::Downloading::%s::TotalSize::%d::Chunks::%d\n", file.DisplayPath, file.Size, len(file.Chunks))
		if err := d.downloadFile(d.Manifest.ID, file); err != nil {
			return err
		}
		fmt.Printf("File::Downloaded::%s::TotalSize::%d::Chunks::%d\n", file.DisplayPath, file.Size, len(file.Chunks))
	}

	return nil
}

func (d *Downloader) DownloadThreaded(size int) error {
	fmt.Println("Download::Start")
	defer fmt.Println("Download::End")

	wait := sync.WaitGroup{}
	limiter := make(chan struct{}, size)

	for _, file := range d.Manifest.Files {
		fmt.Printf("File::Downloading::%s::TotalSize::%d::Chunks::%d\n", file.DisplayPath, file.Size, len(file.Chunks))
		limiter <- struct{}{}
		wait.Add(1)
		go d.downloadThreaded(file, limiter, &wait)
	}

	wait.Wait()

	return nil
}

func (d *Downloader) TestThroughput() error {
	start := time.Now()
	
	resp, err := d.get("http://speedtest.tele2.net/100MB.zip")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}

	duration := time.Since(start).Seconds()
	fmt.Printf("Downloaded::%d::%f::Mbps\n", n, (float64(n) * 8) / (duration * 1_000_000))
	d.maxBytesPerSecond = int64(float64(n) / duration)
	fmt.Printf("MaxBytesPerSecond::%d\n", d.maxBytesPerSecond)

	return nil
}

func (d *Downloader) downloadFile(buildId string, file *chunk.File) error {
	for _, chunk := range file.Chunks {
		if err := d.download(buildId, chunk); err != nil {
			return err
		}

		fmt.Printf("Chunk::Downloaded::%s::TotalSize::%d::Hash::%s\n", file.DisplayPath, chunk.Size, chunk.Hash)
	}

	return nil
}

func (d *Downloader) get(url string) (*http.Response, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *Downloader) download(buildId string, chunk *chunk.Chunk) error {
	exists, err := os.ReadFile(fmt.Sprintf("%s/%s/%s", d.path, buildId, chunk.Hash))
	if err == nil && strings.ToUpper(helpers.MD5(exists)) == chunk.Hash {
		return nil
	}

	resp, err := d.get(fmt.Sprintf("%s/%s/%s", d.root, buildId, chunk.Hash))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download chunk %s for reason %s", chunk.Hash, string(bytes))
	}

	file, err := os.Create(fmt.Sprintf("%s/%s/%s", d.path, buildId, chunk.Hash))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", d.path, buildId), os.ModePerm); err != nil {
			return err	
		}

		file, err = os.Create(fmt.Sprintf("%s/%s/%s", d.path, buildId, chunk.Hash))
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

func (d *Downloader) downloadThreaded(file *chunk.File, limiter chan struct{}, wait *sync.WaitGroup) {
	defer wait.Done()
	defer func() { <-limiter }()

	if file.Check(d.path, d.compilePath) == nil {
		fmt.Printf("File::AlreadyDownloaded::%s::TotalSize::%d\n", file.DisplayPath, file.Size)
		return
	}

	for _, chunk := range file.Chunks {
		if err := d.download(d.Manifest.ID, chunk); err != nil {
			fmt.Printf("Chunk::Failed::%s::TotalSize::%d::Hash::%s::Error::%s\n", file.DisplayPath, chunk.Size, chunk.Hash, err.Error())
			return
		}

		fmt.Printf("Chunk::Downloaded::%s::TotalSize::%d::Hash::%s\n", file.DisplayPath, chunk.Size, chunk.Hash)
	}
}