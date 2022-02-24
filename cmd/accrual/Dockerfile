FROM golang:1.17-alpine3.15
RUN mkdir /app
ADD ./cmd/accrual/accrual_linux_amd64 /app/
WORKDIR /app
CMD ["/app/accrual_linux_amd64"]
