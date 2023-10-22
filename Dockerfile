# a container to create binary in a container for deploy
FROM golang:1.21.3-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# --------------------------------------------------------

# a container for deploy
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

#---------------------------------------------------------

# hot reload environment for local environment
FROM golang:1.21.3 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
