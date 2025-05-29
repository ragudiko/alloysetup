FROM golang:1.24.2 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]
