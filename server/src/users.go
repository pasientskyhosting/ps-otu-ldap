package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// User def
type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	GroupName  string `json:"group_name"`
	ExpireTime int64  `json:"expire_time"`
	CreateTime int64  `json:"create_time"`
	CreateBy   string `json:"create_by"`
}

// UserDB def
type UserDB struct {
	Username   string
	Password   string
	GroupID    int64
	ExpireTime int64
	CreateTime int64
	CreateBy   string
}

// PrepareNewUser - returns new User
func (s *server) PrepareNewUser(userID string, groupName string) (UserDB, error) {

	// Check group exists
	GroupDB, err := s.GetGroup(groupName)
	userDB := UserDB{}

	if err != nil {
		return userDB, fmt.Errorf("Group %s not found: %s", groupName, err)
	}

	userDB.Username = fmt.Sprintf("%s-%s-%s", userID, GroupDB.GroupName, randSeq(8))
	userDB.Password = fmt.Sprintf(randSeq(12))
	userDB.ExpireTime = time.Now().Unix() + int64(GroupDB.LeaseTime)
	userDB.GroupID = GroupDB.id

	return userDB, nil
}

// ForceExpireUsersInGroup - Expire current users in group
func (s *server) ForceExpireUsersInGroup(groupID int64, createBy string) error {

	update, err := s.db.Prepare("UPDATE users SET expire_time = 1 WHERE group_id=$1 AND create_by=$2;")

	if err != nil {
		return errors.New("Could not prepare statement")
	}

	_, err = update.Exec(groupID, createBy)

	if err != nil {
		return errors.New("Could not update db")
	}

	return nil
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Get session user
	var CreateBy, err = s.getUserID(r)

	if err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	var u User

	// Get URL group name
	u.GroupName = chi.URLParam(r, "GroupName")

	UserDB, err := s.PrepareNewUser(CreateBy, u.GroupName)

	if err != nil {
		log.Printf("User error: %+v", err)
		render.Status(r, 404)
		render.JSON(w, r, nil)
		return
	}

	// before insert then expire current users in this group
	err = s.ForceExpireUsersInGroup(UserDB.GroupID, CreateBy)

	if err != nil {

		log.Printf("ERROR: Cannot delete current users")
		render.Status(r, 500)
		render.JSON(w, r, u)
	}

	// insert
	insert, err := s.db.Prepare("INSERT INTO users (username, password, group_id, expire_time, create_time, create_by) values(?,?,?,?,?,?)")

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: preparing insert statement %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	u.Username = UserDB.Username
	u.Password = UserDB.Password
	u.ExpireTime = time.Now().Unix()
	u.CreateTime = time.Now().Unix()
	u.CreateBy = CreateBy

	_, err = insert.Exec(u.Username, u.Password, UserDB.GroupID, u.ExpireTime, u.CreateTime, u.CreateBy)

	if err != nil {
		// handle this error better than this
		log.Printf("ERROR: Inserting into DB  %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, 201)
	render.JSON(w, r, u)

}

func (s *server) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	users := []User{
		{
			Username:   "kj-proxy-sql-fjhfrghrghghr47545",
			Password:   "encrypted_pass1",
			GroupName:  "proxy-sql",
			CreateBy:   "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
		{
			Username:   "kj-rabbitmq-uhefygryg45456",
			Password:   "encrypted_pass2",
			GroupName:  "rabbitmq",
			CreateBy:   "kj",
			ExpireTime: 1554102608,
			CreateTime: 1554102608,
		},
		{
			Username:   "kj-nginx-12344556",
			Password:   "encrypted_pass3",
			GroupName:  "nginx",
			CreateBy:   "kj",
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
