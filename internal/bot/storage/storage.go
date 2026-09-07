package storage

import "fmt"

type User struct {
	ChatID int64  `json:"chat_id"`
	Email  string `json:"email"`
}

type UsersContracts struct {
	ChatID     int64 `json:"user_id"`
	ContractID int64 `json:"document"`
	Active     bool  `json:"active"`
}

type UserContract struct {
	UserID     int64
	ContractID int64
	Active     bool
}

var users []User

func GetAllUsers() []User {
	return users
}

func Update(u []User) {
	users = u
	fmt.Println("count users: ", len(u))
}

func GetUser(email string) *User {
	for _, v := range users {
		if v.Email == email {
			return &v
		}
	}

	return nil
}
