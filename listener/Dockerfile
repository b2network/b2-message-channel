FROM golang:latest as builder
#ARG GOPROXY=https://goproxy.cn
COPY . /src
RUN apt-get update && \
    cd /src/listener && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/build/app cmd/app/main.go
#    &&  GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o /src/build/listener cmd/listener/main.go
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
COPY --from=builder /src/build/app /usr/bin/app
#COPY --from=builder /src/build/listener /usr/bin/listener
COPY --from=builder /src/serv/config/config.yaml /src/config/config.yaml
CMD ["/usr/bin/app"]
