@ECHO OFF
ECHO Compiling for MS Windows ...
set GOOS=windows
set GOARCH=amd64
vgo build .\rajce-get.go
ECHO Compiling for Linux ...
set GOOS=linux
set GOARCH=amd64
vgo build -o rajce-get.bin .\rajce-get.go
ECHO Compiling for MAC OS X ...
set GOOS=darwin
set GOARCH=amd64
vgo build -o rajce-get.app .\rajce-get.go