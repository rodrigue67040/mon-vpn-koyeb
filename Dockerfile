FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /wsproxy

FROM alpine:latest
WORKDIR /app
COPY --from=builder /wsproxy .
ENV PORT=8080
ENV TARGET_ADDR=15.235.60.159:443
EXPOSE 8080
CMD ["/app/wsproxy"]