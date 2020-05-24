exec = main

preparation:
	go get github.com/mattn/go-sqlite3	
	go get github.com/veandco/go-sdl2/{sdl,ttf,img,gfx}

run: preparation
	go run main.go 

build: preparation
	go build -o $(exec) main.go
