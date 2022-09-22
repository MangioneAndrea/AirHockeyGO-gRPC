echo off
echo - protocol buffers -
protoc ./gamepb/gamepb.proto --go_out=plugins=grpc:.
echo - server -
cd server; go build; cd ..
echo - client -
cd client; go build; cd ..