FROM golang:1.24

WORKDIR /app

# Install a specific version of air that we know works
RUN go install github.com/cosmtrek/air@v1.40.4

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
