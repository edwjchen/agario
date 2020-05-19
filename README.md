# agario

## Additional Packages:

go get google.golang.org/grpc

go get golang.org/x/net/context

go get github.com/golang/protobuf/proto

## Build Protocol Buffer:
python3 -m grpc_tools.protoc -I. --python_out=client --grpc_python_out=client protocol.proto

