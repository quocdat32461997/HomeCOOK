user:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=api/protos/userpb \
	--go_out=plugins=grpc:./api/protos/userpb \
	api/protos/userpb/user.proto
	
userproxy:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=api/protos/userpb \
	--grpc-gateway_out=logtostderr=true:./api/protos/userpb \
	api/protos/userpb/user.proto

u: user userproxy


run:
	go run cmd/homecook/main.go