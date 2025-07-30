FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /go-gitlab-user-prune ./cmd/main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /go-gitlab-user-prune /go-gitlab-user-prune
ENTRYPOINT ["/go-gitlab-user-prune"]
