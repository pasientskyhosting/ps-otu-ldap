all: users groups

users:
	export JWT_SECRET=secretfarfar; \
	export ENCRYPTION_KEY=thisisverysecurepassword; \
	export API_KEY=hest; \
	export DB_FILE=/Users/kj/Projects/github/pasientskyhosting/ps-otu-ldap/db/otu.db; \
	go test -run TestUsers;

groups:
	export JWT_SECRET=secretfarfar; \
	export ENCRYPTION_KEY=thisisverysecurepassword; \
	export API_KEY=hest; \
	export DB_FILE=/Users/kj/Projects/github/pasientskyhosting/ps-otu-ldap/db/otu.db; \
	go test -run TestGroups;