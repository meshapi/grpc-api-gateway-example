IMAGE_NAME=grpc-api-gateway-examples:local

docker-build:
	@docker build . -t $(IMAGE_NAME)

docker-run:
	@docker run -v $$(pwd):/project -p4000:4000 -p 40000:40000 --rm -it $(IMAGE_NAME) sh -c 'buf generate && go run .'

docker-shell:
	@docker run -v $$(pwd):/project --rm -it $(IMAGE_NAME) sh

# Install all the necessary protoc-plugins.
install-tools:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-openapiv3@latest
	@go install github.com/meshapi/grpc-api-gateway/codegen/cmd/protoc-gen-grpc-api-gateway@latest

# Generate using buf.
generate:
	@buf generate

# Use only protoc without buf, this requires downloading the grpc api gateway annotation files
# if they are not already downloaded. They will reside inside proto-vendor, a git ignored dir.
generate-protoc: fetch-proto-vendor
	@protoc -I proto-vendor -I . \
		--openapiv3_out=gen \
		chat_service.proto user_service.proto

fetch-proto-vendor:
	@if [ ! -d "proto-vendor" ]; then \
		echo "Downloading grpc_api_gateway_proto.tar.gz..."; \
		curl -L -o grpc_api_gateway_proto.tar.gz \
			https://github.com/meshapi/grpc-api-gateway/releases/download/v0.1.0/grpc_api_gateway_proto.tar.gz; \
		mkdir -p proto-vendor; \
		tar -xzf grpc_api_gateway_proto.tar.gz -C proto-vendor; \
		rm grpc_api_gateway_proto.tar.gz; \
		echo "Downloaded and extracted to proto-vendor."; \
	fi

.PHONY: fetch-proto-vendor
