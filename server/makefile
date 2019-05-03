all: linux64 darwin64

run:
	export JWT_SECRET=secretfarfar; \
	export ENCRYPTION_KEY=thisisverysecurepassword; \
	export API_KEY=hest; \
	export DB_FILE=/Users/kj/Projects/github/pasientskyhosting/ps-otu-ldap/db/otu.db; \
	go run src/main.go src/auth.go src/groups.go src/ping.go src/routes.go src/users.go src/ldap.go src/crypt.go

prep:
	mkdir -p bin

linux64: prep
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap

darwin64: prep
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap