# `acid/chunker`

a utility to improve upload and download speeds for large files by splitting them into smaller chunks.

> note: to download a build, you will need at minimum 2x the size of the file you are downloading. this is because the file is chunked and then reassembled. the chunks are then deleted after reassembly.

## features

- supports both a tui or cli interface
- chunk entire folders into n size chunks and generates a manifest file
- both single & multi-threaded chunking available
- reassemble chunks into original file from manifest from local file
- if download fails, resume from last chunk
- gzip the data after chunking and then un-gzip it after reassembly
- given a src url path, download all files usfing multi-threading
- checks if the file is already downloaded by comparing the hash of the file
- cleans up temporary files after download

## coming soon

make a issue if you want features

![image](https://github.com/user-attachments/assets/ed6708eb-2975-4fe1-80dd-0ae9b82c7b0a)

## usage

`./chunker.exe download <manifest_url> <build_dir>`

- `<manifest_url>` is the url to the manifest file. this is the `.acidmanifest` file that is generated after chunking the build.
- `<build_dir>`is the directory the resulting build will be placed in.

## setup

- `git clone https://github.com/retracmp/chunker.git`
- `cd chunker`
- `go build`

### chunk a build

- `./chunker.exe chunk <build_dir>`

- `<build_dir>` is the directory of the build you want to chunk.

### local file hoster

> this is to emulate a cdn to test the download command

- `./chunker.exe hoster`
- will start a local server on `localhost:80` that serves the files in `./builds`

## wants

- calculate the max download speed and download the chunks efficiently
