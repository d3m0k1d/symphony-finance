FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOCACHE=/root/.cache/go-build
ENV GOTOOLCHAIN=auto
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
	go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/uberproxy/

FROM alpine
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
