
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/argocd-backup-s3 ./cmd/argocd-backup-s3

FROM alpine:latest

RUN apk add --no-cache ca-certificates curl

ARG ARGOCD_VERSION="v2.10.0"

ARG TARGETARCH

RUN set -eux; \
    \
    echo "Downloading Argo CD version ${ARGOCD_VERSION} for ${TARGETARCH}..."; \
    curl -sSL -o argocd-linux-${TARGETARCH} https://github.com/argoproj/argo-cd/releases/download/${ARGOCD_VERSION}/argocd-linux-${TARGETARCH}; \
    \
    install -m 755 argocd-linux-${TARGETARCH} /usr/local/bin/argocd; \
    \
    rm argocd-linux-${TARGETARCH}; \
    \
    echo "Verifying Argo CD installation..."; \
    argocd version --client

COPY --from=builder /app/argocd-backup-s3 /usr/local/bin/argocd-backup-s3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ENTRYPOINT ["/usr/local/bin/argocd-backup-s3"]
