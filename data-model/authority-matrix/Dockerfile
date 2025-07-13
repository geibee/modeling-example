FROM golang:1.23-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]