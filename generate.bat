echo off
echo ---------------------- protocol buffers --------------------------
protoc .\positionpb\position.proto --go_out=plugins=grpc:.
echo -------------------------- server --------------------------------
go build .\server\
echo -------------------------- client --------------------------------
go build .\client\