###################################################################
# Node Builder Stage
###################################################################
FROM node:10-alpine as node_builder
ARG babel_env=production
ENV BABEL_ENV ${babel_env}

RUN apk update --no-cache \
    && apk add --no-cache \
    git=2.24.1-r0 \
    openssh=8.1_p1-r0 \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY gui/ .

RUN npm install && \
    npm run build

###################################################################
# GO Builder Stage
###################################################################
FROM golang:alpine AS go_builder

ARG version=none
ENV GO111MODULE on

RUN apk update \
    && apk add --no-cache \
    git=2.24.1-r0 \
    gcc=9.2.0-r3 \
    g++=9.2.0-r3 \
    upx=3.95-r2 \
    && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/pasientskyhosting/ps-otu-ldap

COPY server/src .

# Get dependencies
RUN go mod init github.com/pasientskyhosting/ps-otu-ldap
RUN go get github.com/go-chi/chi@v3.3.4 && \
    go get github.com/go-chi/chi/middleware@v3.3.4 && \
    go get github.com/go-chi/jwtauth@v4.0.3 && \
	go get github.com/go-chi/render@v1.0.1 && \
    go get github.com/dgrijalva/jwt-go@v3.2.0 && \
    go get github.com/mattn/go-sqlite3@v1.10.0

# Get the rest
RUN go get

# Compile project
RUN go build -ldflags "-s -w -X main.version=${version} -X main.date=$(date '+%Y-%m-%dT%H:%M:%S%z')" -o otu-ldap

###################################################################
# Final Stage
###################################################################
FROM alpine:3.10

# Create WORKDIR
WORKDIR /app

# Copy app binary from the Builder stage image
COPY --from=go_builder /go/src/github.com/pasientskyhosting/ps-otu-ldap/otu-ldap .

COPY --from=node_builder /app/public/ ./html

COPY db/ /data/otu-ldap/

VOLUME ["/data/otu-ldap/"]

# Run the userServer command by default when the container starts.
CMD ["/app/otu-ldap"]
