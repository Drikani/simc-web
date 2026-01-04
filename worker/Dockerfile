FROM simulationcraftorg/simc:latest

WORKDIR /jobs

RUN apk add --no-cache go git bash

COPY main.go /app/
COPY go.mod /app/
COPY go.sum /app/

WORKDIR /app
RUN go mod download
RUN go build -o /usr/local/bin/simc-worker main.go

ENTRYPOINT ["/usr/local/bin/simc-worker"]