GRPC_MESSAGES=rpc.proto


all: grpc proto proto_js compile_ts

grpc:
	protoc --go_out=plugins=grpc:. protobuf/grpc/*.proto

proto:
	protoc --go_out=. protobuf/message.proto

proto_js:
	protoc --js_out=library=vis/myproto_libs,binary:. protobuf/message.proto

compile_ts:
	tsc vis/app.ts
