FROM golang:1.17-alpine3.15
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main ./cmd/gophermart
CMD ["/app/main"]
