package downloader

import (
	"net/url"
	"path"
)

type DownloadOptions struct {
	Manifest        string
	BaseURL         string
	BuildDir        string
	TempDownloadDir string
}

func NewDownloadOptions(manifestURL string, buildDir string) DownloadOptions {
	options := DownloadOptions{
		BuildDir:        buildDir,
		TempDownloadDir: path.Join(buildDir, "TemporaryChunks"),
	}
	parsedURL, err := url.Parse(manifestURL)
	if err != nil {
		return options
	}
	options.Manifest = path.Base(parsedURL.Path)
	parsedURL.Path = path.Dir(parsedURL.Path)
	options.BaseURL = parsedURL.String()

	return options
}