FROM node:20-alpine AS frontend-builder

WORKDIR /web

COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

COPY web/ ./
RUN pnpm run build

FROM golang:1.23.3-alpine AS backend-builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w' -o bot cmd/bot/main.go

FROM alpine:latest

RUN apk --no-cache add --no-scripts ca-certificates tzdata

WORKDIR /app

COPY --from=backend-builder /build/bot .
COPY --from=backend-builder /build/web/dist ./web/dist
COPY --from=backend-builder /build/configs ./configs

RUN mkdir -p /data/configs && chmod -R 777 /data

ENV TZ=Asia/Shanghai

EXPOSE 8080

VOLUME ["/data"]

ENTRYPOINT ["./bot"]
