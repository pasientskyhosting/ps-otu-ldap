###################################################################
# Builder Stage                                                    
###################################################################
FROM golang:alpine AS builder

RUN apk update \
    && apk add git gcc g++ \
    && rm  -rf /tmp/* /var/cache/apk/*

WORKDIR /go/src/github.com/pasientskyhosting/ps-otu-ldap

ADD src/ .

# Get dependencies
RUN go get -d . 

# Compiule project
RUN go build -o otu-ldap

###################################################################
# Final Stage                                                    
###################################################################
FROM alpine

# Create WORKDIR
WORKDIR /app

# Copy app binary from the Builder stage image
COPY --from=builder /go/src/github.com/pasientskyhosting/ps-otu-ldap/otu-ldap .

ADD db/ /data/otu-ldap/

# Document that the service uses port 8080
EXPOSE 8080

VOLUME ["/data/otu-ldap/"]

# Run the userServer command by default when the container starts.
CMD ["/app/otu-ldap"]