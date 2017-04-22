GRPC_MESSAGES=rpc.proto


all: grpc proto proto_js

grpc:
	protoc --go_out=plugins=grpc:. protobuf/${GRPC_MESSAGES}

proto:
	protoc --go_out=. protobuf/message.proto

proto_js:
	protoc --js_out=library=vis/myproto_libs,binary:. protobuf/message.proto
