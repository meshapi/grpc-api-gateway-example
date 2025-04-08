FROM golang:1.24-alpine3.21

RUN apk add --no-cache git curl tar
RUN curl -sSL \
	"https://github.com/bufbuild/buf/releases/download/v1.52.1/buf-$(uname -s)-$(uname -m)" \
	-o /usr/local/bin/buf
RUN chmod u+x /usr/local/bin/buf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-openapiv3@latest
RUN go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-grpc-api-gateway@latest
WORKDIR /project
