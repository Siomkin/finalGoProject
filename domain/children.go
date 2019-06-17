package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Children struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string
	GroupID primitive.ObjectID
	UserID primitive.ObjectID
}

func NewChildren() *Children{
	nc := Children{}
	return &nc
}