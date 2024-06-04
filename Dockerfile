FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git

# Install wkhtmltopdf
RUN apt-get install -y wget wkhtmltopdf
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.jammy_amd64.deb


WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM alpine:3.13
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 2053
ENTRYPOINT /go/bin/web-app --port 2053
