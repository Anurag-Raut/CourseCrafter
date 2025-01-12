FROM golang:1.22.1-alpine

WORKDIR /

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "-ldflags=-w -s", "main.go"]



