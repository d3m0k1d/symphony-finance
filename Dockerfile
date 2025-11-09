FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOCACHE=/root/.cache/go-build
ENV GOTOOLCHAIN=auto
RUN apk add bash
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
	go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 ./do.sh build /app/server


RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@v3.26.0
FROM alpine
WORKDIR /app
RUN apk add bash sqlite
COPY --from=builder /app/server /go/bin/goose .
COPY --from=builder /app/migrations ./migrations
COPY entrypoint.sh .
ENTRYPOINT ["./entrypoint.sh"]
