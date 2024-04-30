FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o main

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
COPY main_page.gohtml .
EXPOSE 8080
CMD ["./main"]
