FROM golang:1.12 AS builder
WORKDIR /go/src/bot
COPY . .
RUN go get .
ENV CGO_ENABLED=0
RUN go build -o app .

FROM alpine
EXPOSE 8000
WORKDIR /app
COPY --from=builder /go/src/bot/ .
CMD ["./app"]
