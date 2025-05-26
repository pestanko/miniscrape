FROM golang:1.24 AS build

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN GOOS="linux" GOARCH="amd64" CGO_ENABLED=0 go build -v -o bin/miniscrape

FROM debian:stable-slim

# Create a non-root user.
RUN adduser -u 1000 --system --group --no-create-home user

# Install curl for healthcheck.
RUN apt update && \
    apt install -y curl && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*

# Change the working directory.
WORKDIR /app

RUN chown -R user:user /app

# Use an unprivileged user.
USER user

COPY --chown=user:user --from=build /app/bin ./bin
COPY --chown=user:user --from=build /app/config ./config

RUN mkdir -p /app/runtime

EXPOSE 8080

ENTRYPOINT ["/app/bin/miniscrape", "serve"]
