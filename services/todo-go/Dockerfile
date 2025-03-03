FROM golang:1.23.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo
EXPOSE 3000

CMD ["/todo"]