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

		error := ErrorForbidden()

		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	// Create group struct
	var g Group

	// Unmarshal body and check for errors
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {

		error := ErrorValidation()
		error.AddMessage("body", "Cannot read body")

		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	err = json.Unmarshal(b, &g)

	if err != nil {
		error := ErrorValidation()
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
		error := ErrorLDAPConn()
		error.AddMessage("ldap", fmt.Sprintf("LDAP backend: %s", err))
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	} else if !lcheck.exists {
		error := ErrorAssetNotFound()
		error.AddMessage("ldap_group_name", "LDAP group does not exist")
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	} else if !lcheck.authorized {
		error := ErrorForbidden()
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
		error := ErrorAssetAlreadyExists()
		error.AddMessage("group_name", "Group already exists")
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	g.CreateTime = time.Now().Unix()
	g.CreateBy = LDAPUser.Username
	cstring, err := json.Marshal(g.CustomProperties) // take user input and create string

	// Check for null and set empty array instead
	if string(cstring) == "null" {
		cstring = []byte("[]")
	}

	insert, err := s.db.Prepare("INSERT INTO groups (group_name, ldap_group_name, lease_time, custom_properties, deleted, create_by, create_time) values(?,?,?,?,?,?,?)")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: preparing insert statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	// Use group name as LDAP group name
	_, err = insert.Exec(g.GroupName, g.LdapGroupName, g.LeaseTime, cstring, 0, g.CreateBy, g.CreateTime)

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

func (g *Group) validateCreateGroup(s *server) ErrorAPI {

	err := ErrorValidation()

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

func (g *Group) validateUpdateGroup(s *server) ErrorAPI {

	err := ErrorValidation()

	if g.GroupName != "" && (len(g.GroupName) < 2 || len(g.GroupName) > 120) {
		err.AddMessage("group_name", "The group_name field must be between 2-120 chars!")
	}

	// check the title field is between 3 to 120 chars
	if g.LeaseTime > 0 && (g.LeaseTime < 60 || g.LeaseTime > 20160) {
		err.AddMessage("lease_time", "The lease_time field must be between 1h (60) and 1y (20160)")
	}

	return err
}

// GetGroup - returns group from db
func (s *server) GetGroup(groupName string) (GroupDB, error) {

	sqlStatement := `SELECT id, ldap_group_name, lease_time, custom_properties, create_time, create_by FROM groups WHERE deleted=0 AND group_name=$1;`
	row := s.db.QueryRow(sqlStatement, groupName)

	var id int64
	var ldapGroupName string
	var leaseTime int
	var customProperties string
	var createTime int64
	var createBy string

	switch err := row.Scan(&id, &ldapGroupName, &leaseTime, &customProperties, &createTime, &createBy); err {
	case nil:
		return GroupDB{
			id:               id,
			GroupName:        groupName,
			LdapGroupName:    ldapGroupName,
			LeaseTime:        leaseTime,
			CustomProperties: customProperties,
			CreateTime:       createTime,
			CreateBy:         createBy,
		}, nil
	}

	return GroupDB{}, errors.New("No group was found")
}

func (s *server) GetAllGroups(w http.ResponseWriter, r *http.Request) {

	var groups []Group
	var c []CustomProperties

	rows, err := s.db.Query("SELECT group_name, ldap_group_name, lease_time, custom_properties, create_time, create_by FROM groups WHERE deleted=0 ORDER BY ldap_group_name, group_name;")

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
		var leaseTime int
		var customProperties string
		var createTime int64
		var createBy string

		err = rows.Scan(&groupName, &ldapGroupName, &leaseTime, &customProperties, &createTime, &createBy)

		if err != nil {
			// handle this error better than this
			log.Printf("ERROR: looping through DB rows: %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, nil)
			return
		}

		err = json.Unmarshal([]byte(customProperties), &c)

		if err != nil {
			log.Printf("Error with custom properties: %+v", c)
			error := ErrorValidation()
			error.AddMessage("body", "Cannot parse JSON from DB")
			render.Status(r, error.StatusCode)
			render.JSON(w, r, error)
			return
		}

		groups = append(groups, Group{
			GroupName:        groupName,
			LdapGroupName:    ldapGroupName,
			LeaseTime:        leaseTime,
			CustomProperties: c,
			CreateTime:       createTime,
			CreateBy:         createBy,
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

func (s *server) UpdateGroup(w http.ResponseWriter, r *http.Request) {

	var LDAPUser, err = s.getLDAPUser(r)

	// Check ldap user
	if !LDAPUser.Admin || err != nil {
		render.Status(r, 403)
		render.JSON(w, r, nil)
		return
	}

	gdb, err := s.GetGroup(chi.URLParam(r, "GroupName"))

	if err != nil {
		error := ErrorAssetNotFound()
		error.AddMessage("asset", "Group was not found")
		render.Status(r, 404)
		render.JSON(w, r, error)
		return
	}

	// Create group struct
	var gpatch Group

	// Unmarshal body and check for errors
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		error := ErrorValidation()
		error.AddMessage("body", "Cannot read body")
		render.Status(r, error.StatusCode)
		render.JSON(w, r, error)
		return
	}

	err = json.Unmarshal(b, &gpatch)

	if validErrs := gpatch.validateUpdateGroup(s); len(validErrs.GetMessages()) > 0 {
		render.Status(r, validErrs.StatusCode)
		render.JSON(w, r, validErrs)
		return
	}

	// set group to deleted
	update, err := s.db.Prepare("UPDATE groups SET group_name=$1, lease_time=$2, custom_properties=$3 WHERE id=$5 AND deleted=0;")

	if err != nil {
		log.Printf("Error: Could not prepare statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	// Update only fields that have changed
	if gpatch.GroupName != "" && gpatch.GroupName != gdb.GroupName {
		gdb.GroupName = gpatch.GroupName
	}

	if gpatch.LeaseTime != 0 && gpatch.LeaseTime != gdb.LeaseTime {
		gdb.LeaseTime = gpatch.LeaseTime
	}

	// Create json string from patched object
	cpstr, err := json.Marshal(gpatch.CustomProperties) // take user input and create string

	if string(cpstr) != "null" && string(cpstr) != gdb.CustomProperties && len(gpatch.CustomProperties) > 0 {
		gdb.CustomProperties = string(cpstr)
	} else if string(cpstr) == "[]" {
		gdb.CustomProperties = "[]"
	}

	_, err = update.Exec(gdb.GroupName, gdb.LeaseTime, gdb.CustomProperties, gdb.id)

	if err != nil {
		log.Printf("Could not update db: %s", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return

	}

	// Create grpup return object
	var gr Group
	var c []CustomProperties

	gr.GroupName = gdb.GroupName
	gr.LdapGroupName = gdb.LdapGroupName
	gr.LeaseTime = gdb.LeaseTime

	err = json.Unmarshal([]byte(gdb.CustomProperties), &c)
	gr.CustomProperties = c

	gr.CreateBy = gdb.CreateBy
	gr.CreateTime = gdb.CreateTime

	render.Status(r, 200)
	render.JSON(w, r, gr)
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

func (s *server) GetAllGroupsInLDAPScope(w http.ResponseWriter, r *http.Request) {

	lg := chi.URLParam(r, "LDAPGroupName")
	var groups []Group
	var c []CustomProperties

	// Get all users in a specific group
	rows, err := s.db.Query("SELECT groups.group_name, groups.ldap_group_name, groups.lease_time, groups.custom_properties, groups.create_time, groups.create_by FROM groups WHERE groups.deleted=0 AND groups.ldap_group_name=$1;", lg)

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
		var customProperties string
		var createTime int64
		var createBy string

		err = rows.Scan(&groupName, &ldapGroupName, &LeaseTime, &customProperties, &createTime, &createBy)

		if err != nil {
			// handle this error better than this
			log.Printf("ERROR: looping through DB rows: %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, nil)
			return
		}

		err = json.Unmarshal([]byte(customProperties), &c)

		if err != nil {
			log.Printf("Error with custom properties: %+v", c)
			error := ErrorValidation()
			error.AddMessage("body", "Cannot parse JSON from DB")
			render.Status(r, error.StatusCode)
			render.JSON(w, r, error)
			return
		}

		groups = append(groups, Group{
			GroupName:        groupName,
			LdapGroupName:    ldapGroupName,
			CustomProperties: c,
			LeaseTime:        LeaseTime,
			CreateTime:       createTime,
			CreateBy:         createBy,
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

	render.Status(r, 200)
	render.JSON(w, r, groups)
	return

}
