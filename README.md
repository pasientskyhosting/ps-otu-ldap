# Patientsky One Time User

## Description
Patientsky One Time User is a small application that runs a Golang API backend and ReactJS frontend. The purpose of this application is to let users login with LDAP and create temporary users that have an expire time.

## Quickstart

### Step 1 - Setup env

```
export ENCRYPTION_KEY=thisisaverysecureencryptionkeyfordatabasepasswords
export LDAP_BIND_DN=uid=bind,cn=sysaccounts,cn=accounts,dc=pasientsky,dc=no
export LDAP_SERVER=fqdn.ldapserver.com
export LDAP_BIND_PASSWORD=thisistheldapbinduserpassword
export JWT_SECRET=jwtsecretforusertokens
export API_KEY=apikeyforothersystemstouse
```

This will use the database located in this repo at `db/otu.db` change it with the env `DB_FILE`

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
`http://localhost:8051/v1/api/`

Check connectivity to API by running:
`curl -X GET http://localhost:8081/v1/api/ping`

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



