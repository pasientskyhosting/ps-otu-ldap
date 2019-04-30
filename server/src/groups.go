package main

import (
	"database/sql"
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

	var userID, err = s.getUserID(r)

	if err != nil {
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

	if validErrs := g.validateCreateGroup(); len(validErrs) > 0 {

		err := map[string]interface{}{"validation_error": validErrs}

		render.Status(r, 400)
		render.JSON(w, r, err)
		return
	}

	g.CreateTime = time.Now().Unix()
	g.CreateBy = userID

	// insert
	insert, err := s.db.Prepare("INSERT INTO groups (group_name, ldap_group_name, lease_time, deleted, create_by, create_time) values(?,?,?,?,?,?)")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: preparing insert statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	_, err = insert.Exec(g.GroupName, g.LdapGroupName, g.LeaseTime, "0", g.CreateBy, g.CreateTime)

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: Inserting into DB  %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, 201)
	render.JSON(w, r, g)

}

func (g *Group) validateCreateGroup() url.Values {

	errs := url.Values{}

	// check if the title empty
	if g.GroupName == "" {
		errs.Add("group_name", "The group_name field is required!")
	}

	// check the title field is between 3 to 120 chars
	if len(g.GroupName) < 2 || len(g.GroupName) > 120 {
		errs.Add("group_name", "The group_name field must be between 2-120 chars!")
	}

	// check if the title empty
	if g.LdapGroupName == "" {
		errs.Add("ldap_group_name", "The ldap_group_name field is required!")
	}

	// check the title field is between 3 to 120 chars
	if len(g.LdapGroupName) < 2 || len(g.LdapGroupName) > 120 {
		errs.Add("ldap_group_name", "The ldap_group_name field must be between 2-120 chars!")
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

	g := chi.URLParam(r, "GroupName")

	users := []User{
		{
			Username:   "kj-" + g + "-fjhfrghrghghr47545",
			Password:   "encrypted_pass1",
			GroupName:  g,
			CreateBy:   "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
		{
			Username:   "ak-" + g + "-fjhfrghrghghr47545",
			Password:   "encrypted_pass1",
			GroupName:  g,
			CreateBy:   "ak",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
	}

	u, _ := json.Marshal(users)

	encryptedPayload := encrypt(u, s.env.ekey)

	render.PlainText(w, r, string(encryptedPayload))
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
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
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

	fmt.Printf("%+v", row)

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

	// g := chi.URLParam(r, "GroupName")

	render.Status(r, 204)
	render.JSON(w, r, nil)
	return

}
