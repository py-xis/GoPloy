# ---------- BUILD STAGE ----------
FROM --platform=linux/arm64 golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary targeting ARM64
RUN go build -o build-server ./cmd/build-server

# ---------- FINAL STAGE ----------
FROM --platform=linux/arm64 ubuntu:22.04

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl git ca-certificates gnupg && \
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean

# Working directory
WORKDIR /home/app

# Copy compiled Go binary
COPY --from=builder /app/build-server .
RUN chmod +x build-server

# Copy script files
COPY main.sh .
RUN chmod +x main.sh

# Entrypoint
ENTRYPOINT ["/home/app/main.sh"]