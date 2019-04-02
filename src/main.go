// lc := ldapConn{
// 	Base:         "cn=accounts,dc=pasientsky,dc=no",
// 	Server:       "odn-glauth01.privatedns.zone",
// 	Port:         636,
// 	BindDN:       "uid=bind,cn=sysaccounts,cn=accounts,dc=pasientsky,dc=no",
// 	BindPassword: bindPass,
// }

package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	_ "github.com/mattn/go-sqlite3"
)

var dbFile = os.Getenv("DB_FILE")                      // SQLite DB path
var port = os.Getenv("PORT")                           // Server listen port
var listen = os.Getenv("LISTEN")                       // Listen
var atoken = os.Getenv("API_TOKEN")                    // API token for sync services
var ekey = os.Getenv("ENCRYPTION_KEY")                 // Encryption key
var jwtSecret = os.Getenv("JWT_SECRET")                // JWT signing key
var ldapBase = os.Getenv("LDAP_BASE")                  // LDAP Base
var ldapServer = os.Getenv("LDAP_SERVER")              // LDAP server url
var ldapBindDN = os.Getenv("LDAP_BIND_DN")             // Bind readonly user
var ldapBindPassword = os.Getenv("LDAP_BIND_PASSWORD") // Bind readonly pass

type server struct {
	db     *sql.DB
	router *chi.Mux
	token  *jwtauth.JWTAuth
	lc     *ldapConn
}

func init() {

	if jwtSecret == "" {
		log.Fatalf("env JWT_SECRET not set!")
	}

	if atoken == "" {
		log.Fatalf("env API_TOKEN not set!")
	}

	if ekey == "" {
		log.Fatalf("env ENCRYPTION_KEY not set!")
	}

	if dbFile == "" {
		dbFile = "/data/otu-ldap/otu.db"
	}

	if port == "" {
		port = "8080"
	}

	if listen == "" {
		listen = "localhost"
	}

	if ldapBase == "" {
		ldapBase = "cn=accounts,dc=pasientsky,dc=no"
	}

	if ldapServer == "" {
		ldapServer = "odn-glauth01.privatedns.zone"
	}

	if ldapBindDN == "" {
		log.Fatalf("env LDAP_BIND_DN not set!")
	}

	if ldapBindPassword == "" {
		log.Fatalf("env LDAP_BIND_PASSWORD not set!")
	}

	// Check dbFile exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Fatalf("Could find db: %q", err)
	}

}

func main() {

	// Setup db conn
	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Could not open db: %q", err)
	}

	lc := ldapConn{
		Base:         ldapBase,
		Server:       ldapServer,
		Port:         636,
		BindDN:       ldapBindDN,
		BindPassword: ldapBindPassword,
	}

	defer db.Close()

	// New server
	s := server{
		router: chi.NewRouter(),
		db:     db,
		token:  jwtauth.New("HS256", []byte(jwtSecret), nil),
		lc:     &lc,
	}

	// For debugging/example purposes, we generate and print a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := s.token.Encode(jwt.MapClaims{"user_id": 123})

	log.Printf("DEBUG: A sample jwt is %s\n", tokenString)

	log.Printf("Started REST API on %s:%s with db %s\n", listen, port, dbFile)
	log.Fatal(http.ListenAndServe(listen+":"+port, s.Routes()))

}
