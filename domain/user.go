package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string
	Pass string
	Role string
	FName string
	LName string
}

func NewUser() *User{
	nu := new(User)
	return nu
}

func (u *User) SetNewName(name string){
	u.Name = name
}

func (u *User) SetRole(role string){
	u.Role = role
}


