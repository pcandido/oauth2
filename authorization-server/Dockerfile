FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY ./private_key.pem ./public_key.pem ./
COPY --from=0 /app/main .
CMD ["./main"]