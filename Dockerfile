FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/argocd-backup-s3 ./cmd/argocd-backup-s3

FROM alpine:latest

RUN apk add --no-cache ca-certificates

ARG TARGETARCH

ARG ARGOCD_VERSION="v2.10.0"

RUN wget https://github.com/argoproj/argo-cd/releases/download/${ARGOCD_VERSION}/argocd-linux-${TARGETARCH} -O /usr/local/bin/argocd

RUN chmod +x /usr/local/bin/argocd

COPY --from=builder /app/argocd-backup-s3 /usr/local/bin/argocd-backup-s3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/usr/local/bin/argocd-backup-s3"]
