package auth

import (
	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/pkg/utils"
)

// LoginCredentials represents a login credentials
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// FindUser instance within the users slice
func FindUser(users []config.User, cred LoginCredentials) *config.User {
	return utils.FindInSlice(users, func(u config.User) bool {
		return u.Username == cred.Username && u.Password == cred.Password
	})
}
