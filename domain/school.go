package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string
}

func NewSchool() *School{
	ns := School{}
	return &ns
}
