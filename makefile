all: stop volume build run

stop: 
	-docker stop otu-ldap

volume:
	docker volume create volume-otu-ldap

run:
	docker run --rm --name=otu-ldap \
	-e ENCRYPTION_KEY=$(ENCRYPTION_KEY) \
	-e DB_FILE=/data/otu-ldap/otu.db \
	-e LDAP_BIND_DN=$(LDAP_BIND_PASSWORD) \
	-e LDAP_SERVER=$(LDAP_BIND_PASSWORD) \
	-e LDAP_BIND_PASSWORD=$(LDAP_BIND_PASSWORD) \
	-e JWT_SECRET=$(JWT_SECRET) \
	-e API_KEY=$(API_KEY) \
	-e LISTEN=0.0.0.0:8081 \
	-v /Users/kj/Projects/github/pasientskyhosting/ps-otu-ldap/db:/data/otu-ldap \
	-p 8081:8081 \
	pasientskyhosting/ps-otu-ldap:latest

build:
	docker build -t pasientskyhosting/ps-otu-ldap:latest .