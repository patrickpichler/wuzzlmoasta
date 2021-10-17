package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var InvalidUsernameOrPassword = errors.New("Invalid username or password")

type UserStore interface {
	// tries to login the given user with the given password.
	// if the login is successfull a token will be generated, which
	// can be used as a cookie to verify a user is logged in.
	// otherwise an error is returned.
	TryLogin(username, password string) (string, error)

	GetUserByName(name string) *ViewableUser

	ValidateToken(token string) (bool, *ViewableUser)

	InvalidateToken(token string) bool
}

type ViewableUser struct {
	Name  string
	Roles []string
}

func BuildInMemoryStore() UserStore {
	return &inMemoryUserStore{
		users: []inMemoryUser{
			{"admin", mustEncrypt("admin"), []string{"admin"}},
			{"seppl", mustEncrypt("1234"), []string{}},
		},
	}
}

func mustEncrypt(password string) []byte {
	ret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return ret
}
