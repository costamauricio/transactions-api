FROM golang:1.20-alpine
RUN apk update && apk upgrade
RUN apk add --no-cache build-base

WORKDIR /app

COPY . .
RUN go mod download

# We need to enable CGO since go-sqlite3 uses it
RUN CGO_ENABLED=1 GOOS=linux go build -o application cmd/api/main.go

# Next stage just copy the binary and runs the application

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite

WORKDIR /app

ENV PORT=80

COPY --from=0 /app/application ./
CMD ["./application"]
