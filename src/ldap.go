package main

import (
	"crypto/tls"
	"errors"
	"fmt"

	"gopkg.in/ldap.v2"
)

type ldapConn struct {
	Base         string
	Server       string
	Port         int
	BindDN       string
	BindPassword string
}

func (lc *ldapConn) LDAPAuthentication(a Auth) error {

	// The username and password we want to check
	username := a.Username
	password := a.Password

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", lc.Server, lc.Port), &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(lc.BindDN, lc.BindPassword)

	if err != nil {
		return err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s)(objectClass=posixaccount))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)

	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return errors.New("User does not exist or too many entries returned")
	}

	if err := l.Bind(sr.Entries[0].DN, password); err != nil {
		return errors.New("Failed to auth: " + username)
	} else {
		fmt.Printf("Authenticated successfuly as %s!\n", username)
		return nil
	}

}
