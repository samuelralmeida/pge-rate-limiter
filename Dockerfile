FROM golang:1.22.10

RUN go install github.com/mitranim/gow@latest

WORKDIR /app

COPY . .

RUN go mod download

CMD [ "gow", "run", "main.go" ]
