FROM golang:1.22.10

RUN go install github.com/mitranim/gow@latest

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8083

CMD [ "gow", "run", "main.go" ]
