APP_NAME=chunker

build:
	go build -ldflags -H=windowsgui -o bin/$(APP_NAME)_gui.exe
	go build -o bin/$(APP_NAME).exe
	echo "build done"