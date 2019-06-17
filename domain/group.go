package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string
	SchoolID primitive.ObjectID
}

func NewGroup() *Group{
	ng := Group{}
	return &ng
}
