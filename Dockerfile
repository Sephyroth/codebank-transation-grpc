FROM golang:1.21-bookworm

WORKDIR /go/src
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential protobuf-compiler librdkafka-dev wget && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install github.com/spf13/cobra/cobra@v1.2.1 && \
    wget https://github.com/ktr0731/evans/releases/download/0.9.1/evans_linux_amd64.tar.gz && \
    tar -xzvf evans_linux_amd64.tar.gz && \
    mv evans ../bin && rm -f evans_linux_amd64.tar.gz && \
    rm -rf /var/lib/apt/lists/*

CMD ["tail", "-f", "/dev/null"]