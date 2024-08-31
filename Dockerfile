FROM golang:1.20-bullseye AS build
RUN apt-get update && apt-get install -y gcc g++ make git

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
wkhtmltopdf \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 2053
ENTRYPOINT /go/bin/web-app --port 2053