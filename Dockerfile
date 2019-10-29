###################################################################
# Node Builder Stage
###################################################################
FROM node:10-alpine as node_builder
ARG babel_env=production
ENV BABEL_ENV ${babel_env}

RUN apk update --no-cache \
    && apk add --no-cache git openssh \
    && rm -rf /var/cache/apk/*

WORKDIR /app

ADD gui/ .

RUN npm install && \
    npm run build

###################################################################
# GO Builder Stage
###################################################################
FROM golang:alpine AS go_builder

ARG version=none

RUN apk update \
    && apk add --no-cache git gcc g++ upx \
    && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/pasientskyhosting/ps-otu-ldap

ADD server/src .

# Get dependencies
RUN go get -d .

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

ADD db/ /data/otu-ldap/

VOLUME ["/data/otu-ldap/"]

# Run the userServer command by default when the container starts.
CMD ["/app/otu-ldap"]
