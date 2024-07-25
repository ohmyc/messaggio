FROM golang

ENV CGO_ENABLE 1

WORKDIR /app

COPY . .

RUN go mod download

RUN go build ./cmd/app