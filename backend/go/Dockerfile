FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/speakup-backend ./cmd/api/main.go

FROM alpine:3.22 AS production

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/speakup-backend ./speakup-backend
COPY --from=build /src/pkg/prompts ./pkg/prompts

EXPOSE 8082

CMD ["./speakup-backend"]
