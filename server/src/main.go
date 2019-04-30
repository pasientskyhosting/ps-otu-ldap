package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
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
	listen           string
	apiKey           string
	ekey             string
	jwtSecret        string
	ldapBase         string
	ldapServer       string
	ldapPort         string
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
	listen string,
	apiKey string,
	ekey string,
	jwtSecret string,
	ldapBase string,
	ldapServer string,
	ldapPort string,
	ldapBindDN string,
	ldapBindPassword string) *env {

	if ldapPort == "" {
		ldapPort = "636"
	}

	if listen == "" {
		listen = "8081"
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
		listen:           listen,
		apiKey:           apiKey,
		ekey:             ekey,
		jwtSecret:        jwtSecret,
		ldapBase:         ldapBase,
		ldapServer:       ldapServer,
		ldapPort:         ldapPort,
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
			os.Getenv("LISTEN"),
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

	log.Printf("Started REST API on localhost:%s with db %s", s.env.listen, s.env.dbFile)

	// Create temp token
	_, ts, _ := s.token.Encode(jwt.MapClaims{"user_id": "apiTest", "exp": jwtauth.ExpireIn(30 * time.Minute)})

	// Start rand seed
	rand.Seed(time.Now().UnixNano())

	log.Printf("Started with temp token: %s", ts)
	log.Fatal(http.ListenAndServe("localhost:"+s.env.listen, s.routes()))

}
