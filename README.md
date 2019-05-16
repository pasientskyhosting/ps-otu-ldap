# PatientSky One Time User

## Description
PatientSky One Time User is a small application that runs a Golang API backend and ReactJS frontend. The purpose of this application is to let users login with LDAP and create temporary users that have an expire time.

## Quickstart

### Step 1 - Setup env

AES only supports key sizes of 16, 24 or 32 bytes. The env `API_ENCRYPTION_KEY` must enfore these limitations.

From doc:
```
// The key argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
```

```
export API_ENCRYPTION_KEY=thisisverysecurepassword && \
export API_LISTEN=0.0.0.0:8081 && \
export API_JWT_SECRET=tokenssecret  && \
export API_KEY=apitsecretkey && \
export API_LDAP_BIND_DN="uid=bind,cn=sysaccounts,cn=accounts,dc=pasientsky,dc=no" && \
export API_LDAP_SERVER=ldap.domain.com && \
export API_LDAP_BIND_PASSWORD=binduserpasswordforldap && \
export API_DB_FILE=/path/to/db.db
```

This will use the database located in this repo at `db/otu.db` change it with the env `API_DB_FILE`

### Step 2 - Build docker image and run

`make all`

You will see similar output:

```
2019/05/03 07:20:23 Server started with env: {dbFile:/data/otu-ldap/otu.db listen:0.0.0.0:8081 apiKey:apikeyforothersystemstouse ekey:thisisaverysecureencryptionkeyfordatabasepasswords jwtSecret:jwtsecretforusertokens ldapBase:cn=accounts,dc=pasientsky,dc=no ldapServer:fqdn.ldapserver.com ldapPort:636 ldapBindDN:uid=bind,cn=sysaccounts,cn=accounts,dc=pasientsky,dc=no ldapBindPassword:thisistheldapbinduserpassword}
2019/05/03 07:20:23 Connected to DB at /data/otu-ldap/otu.db
2019/05/03 07:20:23 Started REST API on 0.0.0.0:8081 with db /data/otu-ldap/otu.db
2019/05/03 07:20:23 Started with temp token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTY4Njk4MjMsInVzZXJfaWQiOiJhcGlUZXN0In0.lLR9rnwLQGIDdsO8q-rywLXPeEfVbpnIIBuvF-FTzQY
``` 

A new temp token will be created each time the server starts for debugging purposes.

### Step 3 - Login

Goto `http://localhost:8081/` and login with your LDAP credentials.

The Golang API will be served by default at:
`http://localhost:8051/api/v1/`

Check connectivity to API by running:
`curl -X GET http://localhost:8081/api/v1/ping`

## HTTPS
We suggest to use a HTTP proxy such as `Traefik` for HTTPS requests infront of the application. 

## Makefile
A makefile exists that will help with the following commands:

### Test
Run self-serving Golang tests with `make test`

### Run
Run frontend and server with `make` or `make all`

### Build
Build image with `make build`

### Stop
Stop docker image from running `make stop`



