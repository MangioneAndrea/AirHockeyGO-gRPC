echo off
echo ---------------------- protocol buffers --------------------------
protoc .\positionpb\position.proto --go_out=plugins=grpc:.
echo -------------------------- server --------------------------------
cd .\server\
go get
go build
cd ..
echo -------------------------- client --------------------------------
cd .\client\
go get
go build
cd ..