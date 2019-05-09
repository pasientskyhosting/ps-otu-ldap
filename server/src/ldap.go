package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"gopkg.in/ldap.v2"
)

type ldapConn struct {
	Base         string
	Server       string
	Port         string
	BindDN       string
	BindPassword string
}

// LDAPUser def
type LDAPUser struct {
	Username    string
	DisplayName string
	Admin       bool
}

// LDAPAuthCheck desc
type LDAPAuthCheck struct {
	authorized bool
	exists     bool
}

// LDAPGroup API response
type LDAPGroup struct {
	LDAPGroupName string `json:"ldap_group_name"`
}

func (s *server) GetAllLDAPGroups(w http.ResponseWriter, r *http.Request) {

	var LDAPUser, err = s.getLDAPUser(r)

	// Check ldap user
	if !LDAPUser.Admin || err != nil {
		render.Status(r, 403)
		render.JSON(w, r, nil)
		return
	}

	// Chech matching LDAP group exits and user has access
	LDAPGroups, err := s.lc.GetAllLDAPGroups()

	if err != nil {
		log.Printf("Error with LDAP backend: %s", LDAPGroups)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: handling DB rows: %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.JSON(w, r, LDAPGroups)
	return
}

// CheckLDAPGroup
func (lc *ldapConn) GetAllLDAPGroups() ([]LDAPGroup, error) {

	port, err := strconv.Atoi(lc.Port)

	if err != nil {
		return nil, errors.New("Error reading ldap port")
	}

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", lc.Server, port), &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return nil, err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(lc.BindDN, lc.BindPassword)

	if err != nil {
		return nil, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprint("(objectClass=posixGroup)"), // The filter to apply
		[]string{"cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	var lg []LDAPGroup

	for _, v := range sr.Entries {
		lg = append(lg, LDAPGroup{LDAPGroupName: strings.Split(strings.Split(v.DN, ",")[0], "cn=")[1]})
	}

	return lg, err

}

// CheckLDAPGroup
func (lc *ldapConn) CheckLDAPGroup(groupName string, username string) (*LDAPAuthCheck, error) {

	lcheck := LDAPAuthCheck{
		authorized: false,
		exists:     false,
	}

	port, err := strconv.Atoi(lc.Port)

	if err != nil {
		return &lcheck, errors.New("Error reading ldap port")
	}

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", lc.Server, port), &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return &lcheck, err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(lc.BindDN, lc.BindPassword)

	if err != nil {
		return &lcheck, err
	}

	// Search for the given groupName
	existsRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s)(objectClass=posixGroup))", groupName), // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	er, err := l.Search(existsRequest)

	if err != nil {
		return &lcheck, err
	}

	if len(er.Entries) != 1 {
		return &lcheck, nil
	}

	// Set exists true
	lcheck.exists = true

	// Search for the given username and get all groups
	accessRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s)(objectClass=posixaccount))", username),
		[]string{"dn", "memberOf", "displayName"},
		nil,
	)

	ar, err := l.Search(accessRequest)

	if len(ar.Entries) != 1 {
		return &lcheck, nil
	}

	for _, v := range ar.Entries[0].Attributes {

		switch v.Name {
		case "memberOf":
			for _, v := range v.Values {
				if v == fmt.Sprintf("cn=%s,cn=groups,cn=accounts,dc=pasientsky,dc=no", groupName) {
					lcheck.authorized = true
				}
			}
			break
		}
	}

	return &lcheck, nil
}

func (lc *ldapConn) LDAPAuthentication(a Auth) (LDAPUser, error) {

	u := LDAPUser{
		Username: a.Username,
		Admin:    false,
	}

	password := a.Password

	port, err := strconv.Atoi(lc.Port)

	if err != nil {
		return LDAPUser{}, errors.New("Error reading ldap port")
	}

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", lc.Server, port), &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return LDAPUser{}, err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(lc.BindDN, lc.BindPassword)

	if err != nil {
		return LDAPUser{}, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s)(objectClass=posixaccount))", u.Username),
		[]string{"dn", "memberOf", "displayName"},
		nil,
	)

	sr, err := l.Search(searchRequest)

	if err != nil {
		return LDAPUser{}, err
	}

	if len(sr.Entries) != 1 {
		return LDAPUser{}, errors.New("User does not exist or too many entries returned")
	}

	accessAllowed := false

	for _, v := range sr.Entries[0].Attributes {
		switch v.Name {
		case "memberOf":
			for _, v := range v.Values {
				if v == "cn=admins,cn=groups,cn=accounts,dc=pasientsky,dc=no" || v == "cn=otu-superheroes,cn=groups,cn=accounts,dc=pasientsky,dc=no" {
					u.Admin = true
					accessAllowed = true
				}
				if v == "cn=otu,cn=groups,cn=accounts,dc=pasientsky,dc=no" {
					accessAllowed = true
				}
			}
			break
		case "displayName":
			u.DisplayName = v.Values[0]
			break
		}
	}

	if err := l.Bind(sr.Entries[0].DN, password); err != nil {
		return LDAPUser{}, errors.New("Failed to auth: " + u.Username)
	}

	if !accessAllowed {
		return LDAPUser{}, errors.New("Not allowed OTU: " + u.Username)
	}

	fmt.Printf("Authenticated successfuly as %s!\n", u.Username)

	return u, nil

}
