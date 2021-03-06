package users

import (
	"time"

	"git.sr.ht/~patrickpichler/wuzzlmoasta/pkg/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type inMemoryUser struct {
	id       uuid.UUID
	username string
	password []byte
	roles    []string
}

type inMemoryUserSession struct {
	username string
	token    string
	created  time.Time
}

type inMemoryUserStore struct {
	users        []inMemoryUser
	userSessions []inMemoryUserSession
}

func (s *inMemoryUserStore) TryLogin(username, password string) (string, error) {
	if s.matchUser(username, password) {
		token := crypto.GenerateRandomString(64)

		s.userSessions = append(s.userSessions, inMemoryUserSession{
			username: username,
			token:    token,
			created:  time.Now(),
		})

		return token, nil
	}

	return "", InvalidUsernameOrPassword
}

func (s *inMemoryUserStore) GetUserByUsername(username string) *ViewableUser {
	for _, u := range s.users {
		if u.username == username {
			ret := new(ViewableUser)
			*ret = ViewableUser{
				Id:    u.id,
				Name:  username,
				Roles: append([]string(nil), u.roles...),
			}

			return ret
		}
	}

	return nil
}

func (s *inMemoryUserStore) ValidateToken(token string) (bool, *ViewableUser) {
	if token == "" {
		return false, nil
	}

	tokenValidStartingWith := time.Now().Add(-1 * 24 * time.Hour)

	for _, us := range s.userSessions {
		if us.token == token {

			if us.created.After(tokenValidStartingWith) {
				return true, s.GetUserByUsername(us.username)
			}

			// token expired
			return false, nil
		}
	}

	return false, nil
}

func (s *inMemoryUserStore) matchUser(username, password string) bool {
	for _, u := range s.users {
		if u.username == username {
			err := bcrypt.CompareHashAndPassword(u.password, []byte(password))

			if err == nil {
				return true
			}
		}
	}

	return false
}

func (s *inMemoryUserStore) InvalidateToken(token string) bool {
	tokenIndex := -1

	for i, s := range s.userSessions {
		if s.token == token {
			tokenIndex = i
			break
		}
	}

	if tokenIndex >= 0 {
		s.userSessions[tokenIndex] = s.userSessions[len(s.userSessions)-1]
		s.userSessions = s.userSessions[:len(s.userSessions)-1]

		return true
	}

	return false
}
