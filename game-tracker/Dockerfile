# Multi-stage build for Game Tracker application

# Stage 1: Build Vue.js Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from stage 1 into cmd/server directory
RUN mkdir -p ./cmd/server/frontend
COPY --from=frontend-builder /app/frontend/dist ./cmd/server/frontend/dist

# Build the server (without dev tag to enable embedding)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/server

# Stage 3: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=backend-builder /app/server .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
