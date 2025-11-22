# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy frontend build from previous stage to expected location for embedding
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o game-tracker ./cmd/server/main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /app/game-tracker .
EXPOSE 8080
CMD ["./game-tracker"]
