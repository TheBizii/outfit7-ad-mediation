# using Go 1.25 bookworm as base image (apparently more stable than alpine)
FROM golang:1.25-bookworm AS base

# move to the /build directory
WORKDIR /build

# copy go.mod and go.sum files to the current directory
COPY go.mod go.sum ./

# install project dependencies
RUN go mod download

# copy the entire project into the container
COPY . .

# build the app
RUN go build -o go-app ./cmd/server

# document the port that may need to be published
EXPOSE 8080

# start the app
CMD ["/build/go-app"]
