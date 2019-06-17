package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prices struct{
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Date int64 //unix date
	Price float64
	GroupID primitive.ObjectID
}

func NewPrices() *Prices{
	np := new(Prices)
	return np
}
