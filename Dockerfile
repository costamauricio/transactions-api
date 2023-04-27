FROM golang:1.20-alpine

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o application cmd/api/main.go

# ----

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=0 /app/application ./
CMD ["./application"]
