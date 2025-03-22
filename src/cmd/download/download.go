package download

import (
	"acid/chunker/src/compile"
	"acid/chunker/src/downloader"
	"acid/chunker/src/helpers"
	"fmt"
	"os"
	"path"
)

func ArgOrPanic(args []string, index int) string {
	if len(args) <= index {
		panic(fmt.Sprintf("Missing argument at index %d", index))
	}
	return args[index]
}

func Start() error {
	timer := helpers.NewPerformanceTimer()
	defer timer.EndTimer()
	
	options := downloader.NewDownloadOptions(
		ArgOrPanic(os.Args[2:], 0), 
		ArgOrPanic(os.Args[2:], 1),
	)

	downloader := downloader.NewDownloader(options.BaseURL, path.Join(options.TempDownloadDir, "chunks"), options.BuildDir)
	if _, err := downloader.FetchManifest(options.Manifest, options.TempDownloadDir); err != nil {
		return err
	}
	
	compiler := compile.NewCompiler(path.Join(options.TempDownloadDir, options.Manifest), path.Join(options.TempDownloadDir, "chunks"), options.BuildDir)
	if err := compiler.Check(); err == nil {
		if err := compiler.Cleanup(options.TempDownloadDir); err != nil {
			return err
		}
		return err
	}

	if err := downloader.DownloadThreaded(4); err != nil {
		return err
	}
	if err := compiler.Compile(); err != nil {
		return err
	}
	if err := compiler.Cleanup(options.TempDownloadDir); err != nil {
		return err
	}

	return nil
}