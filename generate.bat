echo off
go mod vendor
echo - protocol buffers -
protoc .\gamepb\gamepb.proto --go_out=plugins=grpc:.
echo - server -
go build .\server\
echo - client -
set GOARCH=wasm
set GOOS=js
go build -o .\client\dist\main.wasm .\client\main.go

echo - webserver -
set GOARCH=
set GOOS=
go build -o .\client\dist\webserver.exe   .\client\dist\main.go 