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

type CreateWaygateLink struct {
	Name      string
	Link      string
	WaygateId int
}

type WaygateLink struct {
	ID        int
	Name      string
	Link      string
	WaygateId int
}
