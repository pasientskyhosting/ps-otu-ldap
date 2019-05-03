all: stop build run

stop: 
	-docker stop otu-ldap

run:
	docker run --rm --name=otu-ldap \
	-e ENCRYPTION_KEY=$(ENCRYPTION_KEY) \
	-e DB_FILE=/data/otu-ldap/otu.db \
	-e LDAP_BIND_DN=$(LDAP_BIND_DN) \
	-e LDAP_SERVER=$(LDAP_SERVER) \
	-e LDAP_BIND_PASSWORD=$(LDAP_BIND_PASSWORD) \
	-e JWT_SECRET=$(JWT_SECRET) \
	-e API_KEY=$(API_KEY) \
	-e LISTEN=0.0.0.0:8081 \
	-v `pwd`/db:/data/otu-ldap \
	-p 8081:8081 \
	pasientskyhosting/ps-otu-ldap:latest

build:
	docker build -t pasientskyhosting/ps-otu-ldap:latest .

test:
	export ENCRYPTION_KEY=$(ENCRYPTION_KEY); \
	export DB_FILE=`pwd`/db/otu.db; \
	export LDAP_BIND_DN=$(LDAP_BIND_DN); \
	export LDAP_SERVER=$(LDAP_SERVER); \
	export LDAP_BIND_PASSWORD=$(LDAP_BIND_PASSWORD); \
	export JWT_SECRET=$(JWT_SECRET); \
	export API_KEY=$(API_KEY); \
	export LISTEN=0.0.0.0:8081; \
	cd server/src && go test;