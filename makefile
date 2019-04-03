all: stop volume build run

stop: 
	-docker stop otu-ldap

volume:
	docker volume create volume-otu-ldap

run:
	docker run --rm --name=otu-ldap \
	-e ENCRYPTION_KEY=test \
	-e LDAP_BIND_DN="uid=bind,cn=sysaccounts,cn=accounts,dc=pasientsky,dc=no" \
	-e LDAP_BIND_PASSWORD=bindpassword \
	-e LDAP_SERVER=odn-glauth01.privatedns.zone \
	-e JWT_SECRET=jwtsupersecret \
	-e API_TOKEN=hest \
	-e LISTEN=0.0.0.0 \
	-v volume-otu-ldap:/data/otu-ldap \
	-p 8080:8080 \
	pasientskyhosting/ps-otu-ldap:latest

build:
	docker build -t pasientskyhosting/ps-otu-ldap:latest .