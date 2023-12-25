FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go build -o outlier-detection ./cmd/main.go

FROM busybox
WORKDIR /app
COPY --from=builder /app/outlier-detection .

EXPOSE 8080

ENTRYPOINT ["/app/outlier-detection"]
