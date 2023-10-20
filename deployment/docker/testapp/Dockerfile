# Start from golang base image
FROM golang:1.21.1 as builder

# Set the current working directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Config git and Go
RUN export GOSUMDB=off \
	&& go mod download

# Copy sources
COPY . .

# Build the Go app
RUN make build

FROM alpine:3.12

RUN apk add ca-certificates

WORKDIR /app

COPY --from=builder /build/bin/* ./

CMD ["./app"]