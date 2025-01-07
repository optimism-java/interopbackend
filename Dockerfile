FROM golang:latest as builder
ARG GOPROXY=https://goproxy.cn
WORKDIR /src

# Install build dependencies using apt (Debian-based image)
RUN apt-get update && apt-get install -y gcc libc6-dev

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main main.go

# Create final lightweight image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/build/main /usr/bin/main

# Run the binary
CMD ["/usr/bin/main"]
