FROM golang:alpine3.22 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/argocd-backup-s3 ./cmd/argocd-backup-s3

FROM alpine:latest

RUN apk add --no-cache ca-certificates

ARG ARGOCD_VERSION="v2.10.0"
RUN wget https://github.com/argoproj/argo-cd/releases/download/${ARGOCD_VERSION}/argocd-linux-amd64 -O /usr/local/bin/argocd
RUN chmod +x /usr/local/bin/argocd

COPY --from=builder /app/argocd-backup-s3 /usr/local/bin/argocd-backup-s3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/usr/local/bin/argocd-backup-s3"]
