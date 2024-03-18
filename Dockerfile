FROM golang:1.22-alpine AS build

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o bin/miniscrape

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin ./bin
COPY --from=build /app/config ./config

RUN apk add --no-cache curl

EXPOSE 8080

ENTRYPOINT ["/app/bin/miniscrape", "serve"]
