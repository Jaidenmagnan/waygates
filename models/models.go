package models

type CreateUser struct {
	ID       int
	Username string
	Email    string
	Password string
}

type User struct {
	ID       int
	Username string
	Email    string
	Password string
}
