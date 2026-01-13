FROM node:20-alpine AS frontend-builder

WORKDIR /web

COPY wrap-web/package.json wrap-web/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

COPY wrap-web/ ./
RUN pnpm run build

FROM golang:1.23.3-alpine AS backend-builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags '-s -w' -o bot cmd/bot/main.go

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

COPY --from=backend-builder /build/bot .
COPY --from=backend-builder /build/web/dist ./web/dist
COPY --from=backend-builder /build/configs ./configs

ENV TZ=Asia/Shanghai

EXPOSE 8080

VOLUME ["/data"]

ENTRYPOINT ["./bot"]
