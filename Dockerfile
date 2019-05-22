FROM golang:1.12 as build

# set workdir path
WORKDIR /go/src/github.com/arzonus/sitemap

# enable go modules
ENV GO111MODULE=on

# added go.mod and go sum
ADD go.mod .
ADD go.sum .

# download dependecies
RUN go mod download

# added another contains
ADD . .

# build container
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o sitemap cmd/sitemap/main.go

# build alpine container
FROM alpine:latest as production

# copy binary from build image
COPY --from=build /go/src/github.com/arzonus/sitemap/sitemap .

# run app
CMD ["./sitemap"]