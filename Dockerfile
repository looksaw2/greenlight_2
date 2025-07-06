FROM golang:1.24-alpine AS builder
WORKDIR /app 
COPY . .
RUN go mod download
RUN go build -o /bin/server ./cmd/main/main.go


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /bin/server .
EXPOSE 8080
CMD [ "./server" ]