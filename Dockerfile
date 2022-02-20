FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/manojankitha/redis-proxy

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit test
#RUN go test ./...

# oddly required explicit go build flags for go docker
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build the Go app
RUN go build -v -o proxy cmd/main.go

# This container exposes port 9000 to the outside world
EXPOSE 9000

# Run the binary program produced by `go install`
CMD ["./proxy"]
