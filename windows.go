//go:build ignore

// sudo apt-get install gcc-mingw-w64-x86-64
//
// Without console windows but with icon run ./windows-package.sh
//
// Withn console windows
//
//go:generate bash -c "GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOFLAGS=-ldflags=-s go build -tags=opengl"

package main
