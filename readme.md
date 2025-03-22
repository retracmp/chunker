# `acid/chunker`

a better fortnite version downloader

> note: to download a build, you will need at minimum 2x the size of the file you are downloading. this is because the file is chunked and then reassembled. the chunks are then deleted after reassembly.

![image](https://github.com/user-attachments/assets/1d020035-c3ff-4dd4-9e2c-0612d13ba300)

## features

- supports both a tui or a raw cli interface
- chunk entire folders into n size chunks and generates a manifest file
- both single & multi-threaded chunking available
- reassemble chunks into original file from manifest from local file
- if download fails, resume from last chunk
- gzip the data after chunking and then un-gzip it after reassembly
- given a src url path, download all files usfing multi-threading
- checks if the file is already downloaded by comparing the hash of the file
- cleans up temporary files after download

## coming soon

> make a issue if you want more features

- rework the download to either print to the console on each step if a flag is used, or if the tui is being used then update some progress bars for a better user experience
- calculate the max download speed and download the chunks efficiently

## usage

for the basic usage just use 

- `./chunker.exe`

if you want to use the raw cli (mainly for embedding in apps or if you want to use your own manifest files):

> download chunked files using a manifest file from a server url:

`./chunker.exe download <manifest_url> <build_dir>`
- `<manifest_url>` is the url to the manifest file. this is the `.acidmanifest` file that is generated after chunking the build.
- `<build_dir>`is the directory the resulting build will be placed in.

> to chunk the files and generate a manifest file run

`./chunker.exe chunk <build_dir>`
- `<build_dir>` is the directory of the build you want to chunk.


> this is to emulate a server to test the download command

- `./chunker.exe hoster`
- will start a local server on `localhost:80` that serves the files in `./builds`

## setup

- `git clone https://github.com/retracmp/chunker.git`
- `cd chunker`
- `go build`
