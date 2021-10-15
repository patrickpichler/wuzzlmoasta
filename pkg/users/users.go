package users

import (
	"errors"
	"time"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

var InvalidUsernameOrPassword = errors.New("Invalid username or password")

type user struct {
	username string
	password []byte
}

type ViewableUser struct {
	Name string
}

type userSession struct {
	username string
	token    string
	created  time.Time
}

var users = []user{
	{"hansi", mustEncrypt("1234")},
}

var userSessions = []userSession{}

func mustEncrypt(password string) []byte {
	ret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return ret
}

func TryLogin(username, password string) (string, error) {
	if matchUser(username, password) {
		token := crypto.GenerateRandomString(64)

		userSessions = append(userSessions, userSession{
			username: username,
			token:    token,
			created:  time.Now(),
		})

		return token, nil
	}

	return "", InvalidUsernameOrPassword
}

func GetUserByName(name string) *ViewableUser {
	for _, u := range users {
		if u.username == name {
			ret := new(ViewableUser)
			*ret = ViewableUser{
				Name: name,
			}
			return ret
		}
	}

	return nil
}

func IsTokenValid(token string) (bool, *ViewableUser) {
	if token == "" {
		return false, nil
	}

	tokenValidStartingWith := time.Now().Add(-1 * 24 * time.Hour)

	for _, us := range userSessions {
		if us.token == token {

			if us.created.After(tokenValidStartingWith) {
				return true, GetUserByName(us.username)
			}

			// token expired
			return false, nil
		}
	}

	return false, nil
}

func matchUser(username, password string) bool {
	for _, u := range users {
		if u.username == username {
			err := bcrypt.CompareHashAndPassword(u.password, []byte(password))

			if err == nil {
				return true
			}
		}
	}

	return false
}
