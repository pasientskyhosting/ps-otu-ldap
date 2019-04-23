package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db     *sql.DB
	router *chi.Mux
	token  *jwtauth.JWTAuth
	lc     *ldapConn
	env    *env
}

type env struct {
	dbFile           string
	apiKey           string
	ekey             string
	jwtSecret        string
	ldapBase         string
	ldapServer       string
	ldapPort         int64
	ldapBindDN       string
	ldapBindPassword string
}

func newServer(e *env) *server {

	// New server
	s := server{
		router: chi.NewRouter(),
		db:     newDb(*e),
		token:  jwtauth.New("HS256", []byte(e.jwtSecret), nil),
		lc:     newLdapConn(*e),
		env:    e,
	}

	s.router.Routes() // register handlers

	return &s
}

func newLdapConn(e env) *ldapConn {

	lc := ldapConn{
		Base:         e.ldapBase,
		Server:       e.ldapServer,
		Port:         e.ldapPort,
		BindDN:       e.ldapBindDN,
		BindPassword: e.ldapBindPassword,
	}

	return &lc
}

func newEnv(

	dbFile string,
	apiKey string,
	ekey string,
	jwtSecret string,
	ldapBase string,
	ldapServer string,
	ldapPort string,
	ldapBindDN string,
	ldapBindPassword string) *env {

	lport, err := strconv.ParseInt(ldapPort, 10, 64)

	if err != nil {
		log.Printf("Could not parse env LDAP_PORT %s as integer %q. using default 636 \n", string(ldapPort), err)
		lport = 636
	}

	if ldapServer == "" {
		log.Fatalf("Could not parse env LDAP_SERVER %s", ldapServer)
	}

	if ldapBindDN == "" {
		log.Fatalf("Could not parse env LDAP_BIND_DN %s", ldapBindDN)
	}

	if ldapBindPassword == "" {
		log.Fatalf("Could not parse env LDAP_BIND_PASSWORD %s", ldapBindPassword)
	}

	if dbFile == "" {
		dbFile = "/data/otu-ldap/otu.db"
	}

	if ldapBase == "" {
		ldapBase = "cn=accounts,dc=pasientsky,dc=no"
	}

	e := env{
		dbFile:           dbFile,
		apiKey:           apiKey,
		ekey:             ekey,
		jwtSecret:        jwtSecret,
		ldapBase:         ldapBase,
		ldapServer:       ldapServer,
		ldapPort:         lport,
		ldapBindDN:       ldapBindDN,
		ldapBindPassword: ldapBindPassword,
	}

	log.Printf("Server started with env: %+v", e)

	return &e

}

func newDb(e env) *sql.DB {

	// Setup db conn
	db, err := sql.Open("sqlite3", e.dbFile)

	if err != nil {
		log.Fatalf("Could not open db: %q", err)
	}

	log.Printf("Connected to DB at %s", e.dbFile)

	return db
}

func main() {

	s := newServer(
		newEnv(os.Getenv("DB_FILE"),
			os.Getenv("API_KEY"),
			os.Getenv("ENCRYPTION_KEY"),
			os.Getenv("JWT_SECRET"),
			os.Getenv("LDAP_BASE"),
			os.Getenv("LDAP_SERVER"),
			os.Getenv("LDAP_PORT"),
			os.Getenv("LDAP_BIND_DN"),
			os.Getenv("LDAP_BIND_PASSWORD"),
		),
	)

	log.Printf("Started REST API on localhost:8080 with db %s", s.env.dbFile)

	// Create temp token
	_, ts, _ := s.token.Encode(jwt.MapClaims{"user_id": "apiTest", "exp": jwtauth.ExpireIn(30 * time.Minute)})

	log.Printf("Started with temp token: %s", ts)
	log.Fatal(http.ListenAndServe("localhost:8080", s.routes()))

}
