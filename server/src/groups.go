package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Group desc
type Group struct {
	GroupName     string `json:"group_name"`
	LdapGroupName string `json:"ldap_group_name"`
	LeaseTime     int    `json:"lease_time"`
	CreateTime    int    `json:"create_time"`
	CreateBy      string `json:"create_by"`
}

func (s *server) CreateGroup(w http.ResponseWriter, r *http.Request) {

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

	g.CreateTime = 1554102608
	g.CreateBy = "kj"

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
	if g.LeaseTime < 3600 || g.LeaseTime > 84600 {
		errs.Add("lease_time", "The lease_time field must be between 1h (3600) and 1y (86400)")
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

func (s *server) GetAllGroups(w http.ResponseWriter, r *http.Request) {

	groups := []Group{
		{
			GroupName:     "rabbitmq",
			LdapGroupName: "rabbitmq",
			LeaseTime:     3600,
			CreateBy:      "kj",
			CreateTime:    1554102608,
		},
		{
			GroupName:     "proxy-sql",
			LdapGroupName: "proxy-sql",
			LeaseTime:     86400,
			CreateBy:      "kj",
			CreateTime:    1554102608,
		},
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
