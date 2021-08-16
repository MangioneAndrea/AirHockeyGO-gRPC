echo off
go mod vendor
echo - protocol buffers -
protoc .\gamepb\gamepb.proto --go_out=plugins=grpc:.
echo - server -
go build .\server\
echo - client -
go build .\client\