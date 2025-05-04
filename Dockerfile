#### 1. BUILD STAGE ####
FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=/go
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o goph-chat .

#### 2. RUNTIME STAGE ####
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/src/app/goph-chat /usr/local/bin/goph-chat
USER nonroot:nonroot
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/goph-chat"]
CMD []