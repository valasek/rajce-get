@ECHO OFF
ECHO Compiling for MS Windows ...
set GOOS=windows
set GOARCH=amd64
vgo build -o .\binaries\rajce-get.exe .\rajce-get.go
ECHO Compiling for Linux ...
set GOOS=linux
set GOARCH=amd64
vgo build -o .\binaries\rajce-get.bin .\rajce-get.go
ECHO Compiling for MAC OS X ...
set GOOS=darwin
set GOARCH=amd64
vgo build -o .\binaries\rajce-get.app .\rajce-get.go