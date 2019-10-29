all: linux64 darwin64

run:
	export API_JWT_SECRET=secretfarfar; \
	export API_ENCRYPTION_KEY=thisisverysecurepassword; \
	export API_KEY=hest; \
	export API_DB_FILE=`pwd`/db/otu.db; \
	go run src/main.go src/auth.go src/groups.go src/ping.go src/routes.go src/users.go src/ldap.go src/crypt.go src/errors.go

test:
	export API_ENCRYPTION_KEY=$(API_ENCRYPTION_KEY); \
	export API_DB_FILE=`pwd`/db/otu.db; \
	export API_LDAP_BIND_DN=$(API_LDAP_BIND_DN); \
	export API_LDAP_SERVER=$(API_LDAP_SERVER); \
	export API_LDAP_BIND_PASSWORD=$(API_LDAP_BIND_PASSWORD); \
	export API_JWT_SECRET=$(API_JWT_SECRET); \
	export API_KEY=$(API_KEY); \
	export API_LISTEN=0.0.0.0:8081; \
	cd src/ && go test;

prep:
	mkdir -p bin

linux64: prep
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap

darwin64: prep
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap
