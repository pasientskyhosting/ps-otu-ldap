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
func (s *server) PrepareNewUser(createBy string, groupName string) (UserDB, error) {

	// Check group exists
	GroupDB, err := s.GetGroup(groupName)
	userDB := UserDB{}

	if err != nil {
		return userDB, fmt.Errorf("Group %s not found: %s", groupName, err)
	}

	cipherKey := []byte(s.env.ekey)
	password := fmt.Sprintf(randSeq(12))
	encryptedPassword, err := encryptHash(cipherKey, password)

	if err != nil {
		return userDB, fmt.Errorf("Error while encryptHash password: %s", err)
	}

	userDB.Username = fmt.Sprintf("%s-%s-%s", createBy, GroupDB.GroupName, randSeq(8))
	userDB.Password = encryptedPassword
	userDB.ExpireTime = time.Now().Unix() + int64(GroupDB.LeaseTime)
	userDB.GroupID = GroupDB.id
	userDB.CreateTime = time.Now().Unix()
	userDB.CreateBy = createBy

	return userDB, nil
}

// ForceExpireUsersInGroup - Expire current users in group
func (s *server) ExpireUsersInGroup(groupID int64, createBy string) error {

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
	var LDAPUser, err = s.getLDAPUser(r)

	if err != nil {
		render.Status(r, 401)
		render.JSON(w, r, nil)
		return
	}

	var u User

	// Get URL group name
	u.GroupName = chi.URLParam(r, "GroupName")

	UserDB, err := s.PrepareNewUser(LDAPUser.Username, u.GroupName)

	if err != nil {
		log.Printf("User error: %+v", err)
		render.Status(r, 404)
		render.JSON(w, r, nil)
		return
	}

	// before insert then expire current users in this group
	err = s.ExpireUsersInGroup(UserDB.GroupID, LDAPUser.Username)

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
	u.ExpireTime = UserDB.ExpireTime
	u.CreateTime = UserDB.CreateTime
	u.CreateBy = UserDB.CreateBy

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

	var users []User

	// get ssession user
	var LDAPUser, err = s.getLDAPUser(r)

	rows, err := s.db.Query("SELECT users.username, users.password, groups.group_name, users.expire_time, users.create_time, users.create_by FROM users LEFT JOIN GROUPS ON users.group_id = groups.id WHERE groups.deleted=0 and users.create_by=$1 AND users.expire_time > $2;", LDAPUser.Username, time.Now().Unix())

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

	render.JSON(w, r, users)
	return

}
