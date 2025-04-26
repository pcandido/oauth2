package store

import "fmt"

type User struct {
	ID       string
	Email    string
	password string
}

var users = []User{
	{
		ID:       "123",
		Email:    "user@domain.com",
		password: "123",
	},
}

func GetUserByEmail(email string) (*User, error) {
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func GetUserById(id string) (*User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (u *User) ValidatePassword(password string) bool {
	return u.password == password
}
