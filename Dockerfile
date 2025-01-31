FROM golang as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apt update && apt install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/googleapis/gnostic/apps/protoc-gen-openapi

ADD . .

RUN go generate ./...
RUN go build ./...

FROM debian:stable-slim

COPY --from=builder /build/engine /
COPY --from=builder /build/cl /

CMD ["./engine", "--db-in-memory"]
