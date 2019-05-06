package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Group desc
type Group struct {
	GroupName     string `json:"group_name"`
	LdapGroupName string `json:"ldap_group_name"`
	LeaseTime     int    `json:"lease_time"`
	CreateTime    int64  `json:"create_time"`
	CreateBy      string `json:"create_by"`
}

// GroupDB desc
type GroupDB struct {
	id            int64
	GroupName     string
	LdapGroupName string
	LeaseTime     int
	CreateTime    int64
	CreateBy      string
}

func (s *server) CreateGroup(w http.ResponseWriter, r *http.Request) {

	var LDAPUser, err = s.getLDAPUser(r)

	// Check ldap user
	if !LDAPUser.Admin || err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	var g Group

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	err = json.Unmarshal(b, &g)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	if validErrs := g.validateCreateGroup(s); len(validErrs) > 0 {

		err := map[string]interface{}{"validation_error": validErrs}

		render.Status(r, 400)
		render.JSON(w, r, err)
		return
	}

	// Check if already exists
	groupAlreadyExists, _ := s.GetGroup(g.GroupName)

	if groupAlreadyExists.GroupName != "" {
		render.Status(r, 409)
		render.JSON(w, r, nil)
		return
	}

	g.CreateTime = time.Now().Unix()
	g.CreateBy = LDAPUser.Username
	g.LdapGroupName = g.GroupName

	// insert
	insert, err := s.db.Prepare("INSERT INTO groups (group_name, ldap_group_name, lease_time, deleted, create_by, create_time) values(?,?,?,?,?,?)")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: preparing insert statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	// Use group name as LDAP group name
	_, err = insert.Exec(g.GroupName, g.GroupName, g.LeaseTime, "0", g.CreateBy, g.CreateTime)

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: Inserting into DB  %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, 201)
	render.JSON(w, r, g)
	return

}

func (g *Group) validateCreateGroup(s *server) url.Values {

	errs := url.Values{}

	// check if the title empty
	if g.GroupName == "" {
		errs.Add("group_name", "The group_name field is required!")
	}

	if len(g.GroupName) < 2 || len(g.GroupName) > 120 {
		errs.Add("group_name", "The group_name field must be between 2-120 chars!")
	} else {

		// Chech matching LDAP group exits
		exists, err := s.lc.checkLDAPGroupExists(g.GroupName)

		if err != nil {
			errs.Add("group_name", fmt.Sprintf("Error with LDAP backend: %s", err))
		} else if !exists {
			errs.Add("group_name", "LDAP group does not exist with this name")
		}

	}

	// check if the title empty
	if g.LeaseTime == 0 {
		errs.Add("lease_time", "The lease_time field is required!")
	}

	// check the title field is between 3 to 120 chars
	if g.LeaseTime < 60 || g.LeaseTime > 20160 {
		errs.Add("lease_time", "The lease_time field must be between 1h (60) and 1y (20160)")
	}

	return errs
}

func (s *server) GetAllGroupUsers(w http.ResponseWriter, r *http.Request) {

	var users = []User{}
	g, err := s.GetGroup(chi.URLParam(r, "GroupName"))

	if err != nil {
		render.Status(r, 404)
		render.JSON(w, r, nil)
		return
	}

	// Get all users in a specific group
	log.Printf("Get users for group=%d, and expire_time > %d", g.id, time.Now().Unix())
	rows, err := s.db.Query("SELECT users.username, users.password, groups.group_name, users.expire_time, users.create_time, users.create_by FROM users LEFT JOIN GROUPS ON users.group_id = groups.id WHERE groups.deleted=0 AND groups.id=$1 AND users.expire_time > $2;", g.id, time.Now().Unix())

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: connecting to DB: %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	for rows.Next() {

		var username string
		var password string
		var groupName string
		var expireTime int64
		var createTime int64
		var createBy string

		err = rows.Scan(&username, &password, &groupName, &expireTime, &createTime, &createBy)

		if err != nil {
			// handle this error better than this
			log.Printf("ERROR: looping through DB rows: %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, nil)
			return
		}

		users = append(users, User{
			Username:   username,
			Password:   password,
			GroupName:  groupName,
			ExpireTime: expireTime,
			CreateTime: createTime,
			CreateBy:   createBy,
		})
	}

	// get any error encountered during iteration
	err = rows.Err()

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: handling DB rows: %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	u, _ := json.Marshal(users)

	log.Printf("users: %+v, key %s", users, s.env.ekey)

	cipherKey := []byte(s.env.ekey)
	encryptedPayload, err := encryptHash(cipherKey, string(u))

	render.PlainText(w, r, encryptedPayload)
	return

}

// GetGroup - returns group from db
func (s *server) GetGroup(groupName string) (GroupDB, error) {

	sqlStatement := `SELECT id, ldap_group_name, lease_time, create_time, create_by FROM groups WHERE deleted=0 AND group_name=$1;`
	row := s.db.QueryRow(sqlStatement, groupName)

	var id int64
	var ldapGroupName string
	var LeaseTime int
	var createTime int64
	var createBy string

	switch err := row.Scan(&id, &ldapGroupName, &LeaseTime, &createTime, &createBy); err {
	case nil:
		return GroupDB{
			id:            id,
			GroupName:     groupName,
			LdapGroupName: ldapGroupName,
			LeaseTime:     LeaseTime,
			CreateTime:    createTime,
			CreateBy:      createBy,
		}, nil
	}

	return GroupDB{}, errors.New("No group was found")
}

func (s *server) GetAllGroups(w http.ResponseWriter, r *http.Request) {

	var groups []Group

	rows, err := s.db.Query("SELECT group_name, ldap_group_name, lease_time, create_time, create_by FROM groups WHERE deleted=0")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: connecting to DB: %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	for rows.Next() {

		var groupName string
		var ldapGroupName string
		var LeaseTime int
		var createTime int64
		var createBy string

		err = rows.Scan(&groupName, &ldapGroupName, &LeaseTime, &createTime, &createBy)

		if err != nil {
			// handle this error better than this
			log.Printf("ERROR: looping through DB rows: %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, nil)
			return
		}

		groups = append(groups, Group{
			GroupName:     groupName,
			LdapGroupName: ldapGroupName,
			LeaseTime:     LeaseTime,
			CreateTime:    createTime,
			CreateBy:      createBy,
		})
	}

	// get any error encountered during iteration
	err = rows.Err()

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: handling DB rows: %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.JSON(w, r, groups)
	return

}

func (s *server) DeleteGroup(w http.ResponseWriter, r *http.Request) {

	var LDAPUser, err = s.getLDAPUser(r)

	// Check ldap user
	if !LDAPUser.Admin || err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	g, err := s.GetGroup(chi.URLParam(r, "GroupName"))

	if err != nil {
		render.Status(r, 404)
		render.JSON(w, r, nil)
		return
	}

	// expire all users
	err = s.ExpireUsersInGroup(g.id, LDAPUser.Username)

	if err != nil {
		log.Printf("Error: %s", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	// set group to deleted
	update, err := s.db.Prepare("UPDATE groups SET deleted=1 WHERE id=$1;")

	if err != nil {
		log.Printf("Error: Could not prepare statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	_, err = update.Exec(g.id)

	if err != nil {
		log.Printf("Could not update db")
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, 204)
	render.JSON(w, r, nil)
	return

}
