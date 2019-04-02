package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

// User def
type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	GroupName  string `json:"group_name"`
	ExpireTime int    `json:"expire_time"`
	CreateTime int    `json:"create_time"`
	CreatedBy  string `json:"created_by"`
}

func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {

	var u User

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	err = json.Unmarshal(b, &u)

	if err != nil {
		render.Status(r, 400)
		render.JSON(w, r, nil)
		return
	}

	if validErrs := u.validateCreateUser(); len(validErrs) > 0 {

		err := map[string]interface{}{"validation_error": validErrs}

		render.Status(r, 400)
		render.JSON(w, r, err)
		return
	}

	u.Username = "kj-" + u.GroupName
	u.Password = "36746475jhr6hk5"
	u.ExpireTime = 1554102608
	u.CreateTime = 1554102608
	u.CreatedBy = "kj"

	render.JSON(w, r, u)
}

func (u *User) validateCreateUser() url.Values {

	errs := url.Values{}

	// check if the title empty
	if u.GroupName == "" {
		errs.Add("group_name", "The group_name field is required!")
	}

	// check the title field is between 3 to 120 chars
	if len(u.GroupName) < 2 || len(u.GroupName) > 120 {
		errs.Add("group_name", "The group_name field must be between 2-120 chars!")
	}

	return errs
}

func (s *server) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users := []User{
		{
			Username:   "kj-proxy-sql-fjhfrghrghghr47545",
			Password:   "encrypted_pass1",
			GroupName:  "proxy-sql",
			CreatedBy:  "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
		{
			Username:   "kj-rabbitmq-uhefygryg45456",
			Password:   "encrypted_pass2",
			GroupName:  "rabbitmq",
			CreatedBy:  "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
		{
			Username:   "kj-nginx-12344556",
			Password:   "encrypted_pass3",
			GroupName:  "nginx",
			CreatedBy:  "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
	}

	render.JSON(w, r, users)
}

// GetUser def
// func (s *server) GetUser(w http.ResponseWriter, r *http.Request) {

// 	var username string
// 	var password string
// 	var groupID int
// 	var groupName string

// 	username = chi.URLParam(r, "Username")

// 	sqlStatement := "SELECT username, password, group_id, group_name FROM users LEFT JOIN groups ON users.group_id = groups.id WHERE username=$1;"

// 	row := s.db.QueryRow(sqlStatement, username)

// 	switch err := row.Scan(&username, &password, &groupID, &groupName); err {

// 	case sql.ErrNoRows:
// 		log.Printf("No rows were returned for username %v!", username)
// 		render.Status(r, 404)
// 		render.JSON(w, r, nil)
// 		break
// 	default:
// 		render.JSON(w, r, User{
// 			Username:  username,
// 			Password:  password,
// 			GroupID:   groupID,
// 			GroupName: groupName,
// 		})
// 		break
// 	}
// }
