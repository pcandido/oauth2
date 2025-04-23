package store

import "fmt"

type User struct {
	ID       string
	Email    string
	password string
}

var users = map[string]User{
	"user@domain.com": {
		ID:       "123",
		Email:    "user@domain.com",
		password: "123",
	},
}

func GetUser(email string) (*User, error) {
	user, exists := users[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (u *User) ValidatePassword(password string) bool {
	return u.password == password
}
