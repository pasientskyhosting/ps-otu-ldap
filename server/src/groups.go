package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Group desc
type Group struct {
	GroupName        string             `json:"group_name"`
	LdapGroupName    string             `json:"ldap_group_name"`
	CustomProperties []CustomProperties `json:"custom_properties"`
	LeaseTime        int                `json:"lease_time"`
	CreateTime       int64              `json:"create_time"`
	CreateBy         string             `json:"create_by"`
}

// CustomProperties desc
type CustomProperties struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GroupDB desc
type GroupDB struct {
	id               int64
	GroupName        string
	LdapGroupName    string
	CustomProperties string
	LeaseTime        int
	CreateTime       int64
	CreateBy         string
}

func (s *server) CreateGroup(w http.ResponseWriter, r *http.Request) {

	// Get ldap group name from URL
	lg := chi.URLParam(r, "LDAPGroupName")

	var LDAPUser, err = s.getLDAPUser(r)

	// Check ldap user
	if !LDAPUser.Admin || err != nil {

		error := NewForbiddenError()

		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	// Create group struct
	var g Group

	// Unmarshal body and check for errors
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {

		error := NewValidationError()
		error.AddMessage("body", "Cannot read body")

		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	err = json.Unmarshal(b, &g)

	if err != nil {
		error := NewValidationError()
		error.AddMessage("body", "Cannot parse JSON")

		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	// Set LDAP group name
	g.LdapGroupName = lg

	// Chech matching LDAP group exits and user has access
	lcheck, err := s.lc.CheckLDAPGroup(g.LdapGroupName, LDAPUser.Username)

	if err != nil {
		error := LDAPConnError()
		error.AddMessage("ldap", fmt.Sprintf("LDAP backend: %s", err))
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	} else if !lcheck.exists {
		error := NewAssetNotFoundError()
		error.AddMessage("ldap_group_name", "LDAP group does not exist")
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	} else if !lcheck.authorized {
		error := NewForbiddenError()
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	if validErrs := g.validateCreateGroup(s); len(validErrs.GetMessages()) > 0 {
		render.Status(r, validErrs.StatusCode)
		render.JSON(w, r, validErrs)
		return
	}

	// Check if already exists
	groupAlreadyExists, _ := s.GetGroup(g.GroupName)

	if groupAlreadyExists.GroupName != "" {
		error := NewAssetAlreadyExistsError()
		error.AddMessage("group_name", "Group already exists")
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	g.CreateTime = time.Now().Unix()
	g.CreateBy = LDAPUser.Username

	cp, err := json.Marshal(g.CustomProperties)

	if err != nil {
		render.Status(r, 500)
		log.Printf("ERROR: Marshal custom properties %+v", err)
		render.JSON(w, r, nil)
		return
	}

	// insert
	insert, err := s.db.Prepare("INSERT INTO groups (group_name, ldap_group_name, lease_time, custom_properties, deleted, create_by, create_time) values(?,?,?,?,?,?,?)")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: preparing insert statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	// Use group name as LDAP group name
	_, err = insert.Exec(g.GroupName, g.LdapGroupName, g.LeaseTime, cp, 0, g.CreateBy, g.CreateTime)

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

func (g *Group) validateCreateGroup(s *server) APIError {

	err := NewValidationError()

	// check if the title empty
	if g.GroupName == "" {
		err.AddMessage("group_name", "The group_name field is required!")
	}

	if len(g.GroupName) < 2 || len(g.GroupName) > 120 {
		err.AddMessage("group_name", "The group_name field must be between 2-120 chars!")
	}

	// check if the title empty
	if g.LeaseTime == 0 {
		err.AddMessage("lease_time", "The lease_time field is required!")
	}

	// check the title field is between 3 to 120 chars
	if g.LeaseTime < 60 || g.LeaseTime > 20160 {
		err.AddMessage("lease_time", "The lease_time field must be between 1h (60) and 1y (20160)")
	}

	return err
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

	rows, err := s.db.Query("SELECT group_name, ldap_group_name, lease_time, create_time, create_by FROM groups WHERE deleted=0 ORDER BY ldap_group_name, group_name;")

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
		render.Status(r, 403)
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
