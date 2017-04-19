protoc --js_out=library=vis/myproto_libs,binary:. protobuf/*.proto
protoc --go_out=. protobuf/*.proto
