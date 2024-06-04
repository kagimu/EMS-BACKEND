FROM alpine:3.13 AS build

# Install wkhtmltopdf and build dependencies
RUN apk --no-cache add gcc g++ make git wkhtmltopdf

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

# Final image
FROM alpine:3.13
RUN apk --no-cache add ca-certificates

# Copy application binary and wkhtmltopdf from build stage
COPY --from=build /go/src/app/bin/web-app /usr/bin/web-app
COPY --from=build /usr/local/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf  # Adjust path if wkhtmltopdf installs elsewhere

EXPOSE 2053
ENTRYPOINT ["/usr/bin/web-app", "--port", "2053"]