all: stop volume build run

stop: 
	docker stop otu-ldap

volume:
	docker volume create volume-otu-ldap

run:
	docker run --rm --name=otu-ldap \
	-e ENCRYPTION_KEY=test \
	-e LDAP_BIND_DN=$LDAP_BIND_DN \
	-e LDAP_BIND_PASSWORD=$LDAP_BIND_PASSWORD \
	-e JWT_SECRET=jwtsupersecret \
	-e API_TOKEN=hest \
	-e LISTEN=0.0.0.0 \
	-v volume-otu-ldap:/data/otu-ldap \
	-p 8080:8080 \
	pasientskyhosting/ps-otu-ldap:latest

build:
	docker build -t pasientskyhosting/ps-otu-ldap:latest .