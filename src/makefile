all: linux64 darwin64

run:
	export JWT_SECRET=secretfarfar; \
	export ENCRYPTION_KEY=test; \
	export API_TOKEN=hest; \
	export DB_FILE=/Users/kj/Projects/github/pasientskyhosting/ps-otu-ldap/db/otu.db; \
	go run main.go auth.go groups.go ping.go routes.go users.go ldap.go

prep:
	mkdir -p bin

linux64: prep
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap

darwin64: prep
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/otu-ldap