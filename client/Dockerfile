FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/main .
CMD ["./main"]