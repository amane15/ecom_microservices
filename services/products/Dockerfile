FROM golang:1.25.6alpine3.23

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./
COPY ../../proto ./

RUN go mod download

CMD ["air"]
