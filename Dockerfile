# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

# Accept build args for Firebase config
ARG VITE_FIREBASE_API_KEY
ARG VITE_FIREBASE_AUTH_DOMAIN
ARG VITE_FIREBASE_PROJECT_ID
ARG VITE_FIREBASE_STORAGE_BUCKET
ARG VITE_FIREBASE_MESSAGING_SENDER_ID
ARG VITE_FIREBASE_APP_ID
ARG VITE_API_URL

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./

# Set env vars for Vite build
ENV VITE_FIREBASE_API_KEY=$VITE_FIREBASE_API_KEY
ENV VITE_FIREBASE_AUTH_DOMAIN=$VITE_FIREBASE_AUTH_DOMAIN
ENV VITE_FIREBASE_PROJECT_ID=$VITE_FIREBASE_PROJECT_ID
ENV VITE_FIREBASE_STORAGE_BUCKET=$VITE_FIREBASE_STORAGE_BUCKET
ENV VITE_FIREBASE_MESSAGING_SENDER_ID=$VITE_FIREBASE_MESSAGING_SENDER_ID
ENV VITE_FIREBASE_APP_ID=$VITE_FIREBASE_APP_ID
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Copy frontend build from previous stage into cmd/server for embedding
COPY --from=frontend-builder /app/frontend ./cmd/server/frontend

# Build the server (frontend is embedded via go:embed)
RUN cd cmd/server && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../server .

# Stage 3: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy backend binary (frontend is embedded inside)
COPY --from=backend-builder /app/server .

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
