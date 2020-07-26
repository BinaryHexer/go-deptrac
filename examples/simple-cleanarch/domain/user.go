package domain

type User struct {
	Id   string
	Name string
}

type UserRepo interface {
	GetUser(id string) User
	CreateUser(name string) User
}
