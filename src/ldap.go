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
	Port         int64
	BindDN       string
	BindPassword string
}

type ldapUser struct {
	username    string
	displayName string
	admin       bool
}

func (lc *ldapConn) LDAPAuthentication(a Auth) (ldapUser, error) {

	u := ldapUser{
		username: a.Username,
		admin:    false,
	}

	password := a.Password

	l, err := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", lc.Server, lc.Port), &tls.Config{InsecureSkipVerify: true})

	if err != nil {
		return ldapUser{}, err
	}

	defer l.Close()

	// First bind with a read only user
	err = l.Bind(lc.BindDN, lc.BindPassword)

	if err != nil {
		return ldapUser{}, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(uid=%s)(objectClass=posixaccount))", u.username),
		[]string{"dn", "memberOf", "displayName"},
		nil,
	)

	sr, err := l.Search(searchRequest)

	if err != nil {
		return ldapUser{}, err
	}

	if len(sr.Entries) != 1 {
		return ldapUser{}, errors.New("User does not exist or too many entries returned")
	}

	for _, v := range sr.Entries[0].Attributes {
		switch v.Name {
		case "memberOf":
			for _, v := range v.Values {
				if v == "cn=admins,cn=groups,cn=accounts,dc=pasientsky,dc=no" || v == "cn=otu-superheroes,cn=groups,cn=accounts,dc=pasientsky,dc=no" {
					u.admin = true
				}
			}
			break
		case "displayName":
			u.displayName = v.Values[0]
			break
		}
	}

	if err := l.Bind(sr.Entries[0].DN, password); err != nil {
		return ldapUser{}, errors.New("Failed to auth: " + u.username)
	}

	fmt.Printf("Authenticated successfuly as %s!\n", u.username)

	return u, nil

}
