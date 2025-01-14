PROTO_ROOT=./proto


all: link clean build unlink
link: 
	ln -s ~/proto ./proto
unlink:
	unlink ./proto
build:
	protoc -I${PROTO_ROOT} \
               --go_out . \
               --go-grpc_out . \
               --go-grpc_opt=require_unimplemented_servers=false \
               --grpc-gateway_out=logtostderr=true:. \
               ${PROTO_ROOT}/recon/jobmgmt/messages/*.proto 
	protoc -I${PROTO_ROOT} \
               --go_out . \
               --go-grpc_out . \
               --go-grpc_opt=require_unimplemented_servers=false \
               --grpc-gateway_out=logtostderr=true:. \
               ${PROTO_ROOT}/recon/jobmgmt/adminapi/*.proto 
	protoc -I${PROTO_ROOT} \
               --go_out . \
               --go-grpc_out . \
               --go-grpc_opt=require_unimplemented_servers=false \
               --grpc-gateway_out=logtostderr=true:. \
               ${PROTO_ROOT}/recon/jobmgmt/*.proto 

install_gen:
	go install  "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway" \
        "google.golang.org/grpc/cmd/protoc-gen-go-grpc" \
        "google.golang.org/protobuf/cmd/protoc-gen-go" \
        "google.golang.org/protobuf/cmd/protoc-gen-go" 
   
clean:
	rm -fr reconpb/*


