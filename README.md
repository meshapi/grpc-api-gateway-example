# grpc-api-gateway-example

This repository provides examples demonstrating the capabilities of the gRPC API Gateway project.

## Installation

You can either install all the required tools locally to run the project containing the examples or use a Docker container to quickly try things out without setting up a local environment.

### Docker Approach

To build and run the Docker container, execute:

```sh
make docker-build && make docker-run
```

### Local Approach

If you haven't installed `protoc` yet, follow the installation instructions [here](https://protobuf.dev/installation/).

To install the necessary `protoc` plug-ins, run:

```sh
make install-tools
```

[Buf](https://buf.build/docs/cli/installation/) is a powerful tool for protocol buffer code generation. If you have Buf installed or wish to install and use it, you can generate the required files with:

```sh
make generate
```

Alternatively, you can use `protoc` directly. Note that this approach will download some required proto files:

```sh
make generate-protoc
```

## Examples

### Chat Application with Gateway Configurations

This example includes a UI accessible at [http://localhost:4000/web](http://localhost:4000/web). You can open multiple tabs to simulate multiple users interacting with the application. It serves as a simple demonstration of long-lived WebSocket integration.

Refer to the following files for this example:
- `chat_service.go`
- `chat_service.proto`
- `chat_service_gateway.yaml`

### Unary User Endpoint with Server-Sent Events (SSE)

This example showcases a unary user endpoint and an SSE-based user event stream.

Relevant files for this example:
- `user_service.go`
- `user_service.proto`

You can test this example using the following cURL commands:

1. Test the SSE user event stream endpoint, you should be receiving the updates as they come:

```sh
curl -N http://localhost:4000/users-stream
```

2. Add a new user:

```sh
curl -X POST http://localhost:4000/users --data '{"name": "Something"}'
```

3. Delete the user:

```sh
curl -X DELETE http://localhost:4000/users/:id
```
