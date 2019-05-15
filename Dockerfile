FROM golang
WORKDIR /go/src/github.com/arzonus/sitemap
ENV GO111MODULE=on
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN go build -o bin cmd/cli/main.go
CMD ["./bin"]