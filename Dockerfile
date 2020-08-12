###################################################################
# Node Builder Stage
###################################################################
FROM node:10-alpine as node_builder
ARG babel_env=production
ENV BABEL_ENV ${babel_env}

RUN apk update --no-cache \
    && apk add --no-cache \
    git=2.24.3-r0 \
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
    git=2.26.2-r0 \
    gcc=9.3.0-r2 \
    g++=9.3.0-r2 \
    upx=3.96-r0 \
    && rm -rf /var/cache/apk/*

WORKDIR /go/src/github.com/pasientskyhosting/ps-otu-ldap

COPY server/src .

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
