###################################################################
# Node Builder Stage                                                    
###################################################################
FROM node:10-alpine as node_builder
ENV BABEL_ENV production

RUN apk update \
    && apk add git openssh \
    && rm  -rf /tmp/* /var/cache/apk/*

WORKDIR /app

ADD gui/ .

RUN npm install && \
    npm run build

###################################################################
# GO Builder Stage                                                    
###################################################################
FROM golang:alpine AS go_builder

RUN apk update \
    && apk add git gcc g++ upx \
    && rm  -rf /tmp/* /var/cache/apk/*

WORKDIR /go/src/github.com/pasientskyhosting/ps-otu-ldap

ADD server/src .

# Get dependencies
RUN go get -d . 

# Compiule project
RUN go build -ldflags="-s -w" -o otu-ldap

###################################################################
# Final Stage                                                    
###################################################################
FROM alpine

# Create WORKDIR
WORKDIR /app

# Copy app binary from the Builder stage image
COPY --from=go_builder /go/src/github.com/pasientskyhosting/ps-otu-ldap/otu-ldap .

COPY --from=node_builder /app/public/ ./html 

ADD db/otu.db /data/otu-ldap/otu.db

VOLUME ["/data/otu-ldap/"]

# Run the userServer command by default when the container starts.
CMD ["/app/otu-ldap"]