# agario

This project is great :))))))))

## Add repository to GOPATH

export GOPATH=<path_to_repo>

export PATH=$PATH:/$GOPATH/bin


## Additional Packages:

go get google.golang.org/grpc

go get golang.org/x/net/context

go get github.com/golang/protobuf/proto

go get golang.org/x/sys/unix

## Build Protocol Buffer:

Server: 

protoc --go_out=plugins=grpc:server/blob blob.proto

Client: 

python3 -m grpc_tools.protoc -I. --python_out=client --grpc_python_out=client blob.proto

