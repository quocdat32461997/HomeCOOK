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

runu:
	go run cmd/homecook/userssvc/main.go

runc:
	go run cmd/homecook/chefssvc/main.go

usersd: 
	docker build -t gcr.io/homecook/homecook-users:latest cmd/homecook/userssvc/
	docker push gcr.io/homecook/homecook-users:latest
	
chefsd: 
	docker build -t gcr.io/homecook/homecook-chefs:latest cmd/homecook/chefssvc/
	docker push gcr.io/homecook/homecook-chefs:latest

kubeinit:
	gcloud container clusters create homecook
	gcloud codntainer clusters get-credentials homecook
	kubectl create -f kubernetes/deployments/chefs.yaml
	kubectl create -f kubernetes/deployments/users.yaml
	kubectl create -f kubernetes/ingress.yaml

buildu:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o cmd/homecook/userssvc/users cmd/homecook/userssvc/main.go

buildc:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o cmd/homecook/chefssvc/chefs cmd/homecook/chefssvc/main.go

expose:
	kubectl expose deployment homecook-users --target-port=8080 --type=NodePort
	kubectl expose deployment homecook-chefs --target-port=8081 --type=NodePort
