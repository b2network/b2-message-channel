FROM golang:latest as builder
ARG GOPROXY=https://goproxy.cn
COPY . /src
RUN apt-get update && \
    cd /src && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main cmd/app/main.go
# RUN apt-get update && apt install -y protobuf-compiler git && \
#    cd /tmp && git clone https://github.com/googleapis/googleapis.git && \
#    cp -r /tmp/googleapis/* /usr/local/include/ && \
#    cd /src && \
#    go get github.com/golang/protobuf/protoc-gen-go && \
#    go install github.com/golang/protobuf/protoc-gen-go && \
#    PROTO_INCLUDE=/usr/local/include make proto && \
#    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main cmd/app/main.go


# =====
FROM alpine:latest
WORKDIR "/src"
COPY --from=builder /src/build/main /usr/bin/main
COPY --from=builder /src/config/config.yaml /src/config/config.yaml
CMD ["/usr/bin/main"]
