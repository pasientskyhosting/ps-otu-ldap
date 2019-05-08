all: stop build run

stop: 
	-docker stop otu-ldap

run:
	docker run --rm --name=otu-ldap \
	-e API_ENCRYPTION_KEY=$(API_ENCRYPTION_KEY) \
	-e API_DB_FILE=/data/otu-ldap/otu.db \
	-e API_LDAP_BIND_DN=$(API_LDAP_BIND_DN) \
	-e API_LDAP_SERVER=$(API_LDAP_SERVER) \
	-e API_LDAP_BIND_PASSWORD=$(API_LDAP_BIND_PASSWORD) \
	-e API_JWT_SECRET=$(API_JWT_SECRET) \
	-e API_KEY=$(API_KEY) \
	-e API_LISTEN=0.0.0.0:8081 \
	-v `pwd`/db:/data/otu-ldap \
	-p 8081:8081 \
	pasientskyhosting/ps-otu-ldap:latest

build:
	docker build -t pasientskyhosting/ps-otu-ldap:latest .

test:
	export API_ENCRYPTION_KEY=$(API_ENCRYPTION_KEY); \
	export API_DB_FILE=`pwd`/db/otu.db; \
	export API_LDAP_BIND_DN=$(API_LDAP_BIND_DN); \
	export API_LDAP_SERVER=$(API_LDAP_SERVER); \
	export API_LDAP_BIND_PASSWORD=$(API_LDAP_BIND_PASSWORD); \
	export API_JWT_SECRET=$(API_JWT_SECRET); \
	export API_KEY=$(API_KEY); \
	export API_LISTEN=0.0.0.0:8081; \
	cd server/src && go test;