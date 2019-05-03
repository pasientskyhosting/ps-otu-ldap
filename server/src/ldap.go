package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

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
