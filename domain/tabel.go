package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TabelRecord struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Date int64
	Value bool	`json:"value"`
	ChildID primitive.ObjectID
}

func NewTabelRecord() *TabelRecord{
	ntr := new(TabelRecord)
	return ntr
}


