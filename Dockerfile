FROM golang:1.26.6-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /vibrox-arena ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache wget && addgroup -S vibrox && adduser -S -G vibrox vibrox
COPY --from=builder /vibrox-arena /usr/local/bin/vibrox-arena
USER vibrox
EXPOSE 8100 8054
ENTRYPOINT ["/usr/local/bin/vibrox-arena"]
