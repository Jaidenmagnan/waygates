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

type CreateWaygate struct {
	Name   string
	UserId int
}

type Waygate struct {
	ID     int
	Name   string
	UserId int
}
