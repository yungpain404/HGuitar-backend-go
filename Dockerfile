FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
# Dùng image golang chuẩn thay vì dotnet (sửa ở bên dưới)
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /server /server
EXPOSE 5501
ENV SERVER_PORT=5501
ENTRYPOINT ["/server"]
