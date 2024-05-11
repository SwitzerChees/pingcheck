# Use the official Go image to build the binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.22.3 AS builder

# Set the working directory inside the container
WORKDIR /build

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
# -o flag sets the output file name
# CGO_ENABLED=0 disables CGO for a fully static binary
# GOOS=linux sets the target OS to Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /pingcheck .

# Use a minimal image to run the application
FROM alpine:latest

WORKDIR /app

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary to the production image from the builder stage.
COPY --from=builder /pingcheck /app/

# Run the web service on container startup.
CMD ["./pingcheck"]
