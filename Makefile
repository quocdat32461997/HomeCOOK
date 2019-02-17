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

shu:
	sh scripts/user.sh

chef:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=api/protos/chefpb \
	--go_out=plugins=grpc:./api/protos/chefpb \
	api/protos/chefpb/chef.proto
	
chefproxy:
	protoc -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=api/protos/chefpb \
	--grpc-gateway_out=logtostderr=true:./api/protos/chefpb \
	api/protos/chefpb/chef.proto

c: chef chefproxy

shc:
	sh scripts/chef.sh

run:
	go run cmd/homecook/main.go