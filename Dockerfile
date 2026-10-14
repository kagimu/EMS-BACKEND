FROM golang:1.20-bullseye AS build

# Debian Bullseye is end-of-life: its packages have moved from deb.debian.org
# to archive.debian.org, so a plain `apt-get update` now 404s. Point apt at the
# archive (and skip the archive's absent Valid-Until checks) so the build stage
# can still install its tools. The base image and Go version are unchanged.
RUN set -eux; \
    if [ -f /etc/apt/sources.list ]; then \
        sed -i \
            -e 's|deb.debian.org/debian-security|archive.debian.org/debian-security|g' \
            -e 's|security.debian.org/debian-security|archive.debian.org/debian-security|g' \
            -e 's|deb.debian.org/debian|archive.debian.org/debian|g' \
            -e '/bullseye-updates/d' \
            /etc/apt/sources.list; \
    fi; \
    echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/99no-check-valid-until; \
    apt-get update; \
    apt-get install -y gcc g++ make git

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/web-app ./main.go

FROM debian:bullseye-slim

# Same end-of-life remediation for the runtime stage, then install wkhtmltopdf.
# wkhtmltopdf is REQUIRED at runtime: domain/pdfGenerator uses go-wkhtmltopdf,
# which executes the `wkhtmltopdf` binary to render PDFs (used by the statistics
# / report endpoints). Do not remove it.
RUN set -eux; \
    if [ -f /etc/apt/sources.list ]; then \
        sed -i \
            -e 's|deb.debian.org/debian-security|archive.debian.org/debian-security|g' \
            -e 's|security.debian.org/debian-security|archive.debian.org/debian-security|g' \
            -e 's|deb.debian.org/debian|archive.debian.org/debian|g' \
            -e '/bullseye-updates/d' \
            /etc/apt/sources.list; \
    fi; \
    echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/99no-check-valid-until; \
    apt-get update; \
    apt-get install -y wkhtmltopdf; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*

ENV XDG_RUNTIME_DIR=/tmp/runtime-root
RUN mkdir -p /tmp/runtime-root && chmod 700 /tmp/runtime-root


WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 2053
ENTRYPOINT /go/bin/web-app --port 2053
