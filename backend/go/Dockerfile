FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/cosmtrek/air@v1.40.4

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
