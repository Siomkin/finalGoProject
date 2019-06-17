package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main/domain"
)

type PricesRepository interface {
	GetPrice(ctx context.Context, date int64, groupId primitive.ObjectID) (*domain.Prices, error)
	AddPrice(ctx context.Context, date int64, groupId primitive.ObjectID, price float64)
	ChangePrice(ctx context.Context, date int64, groupId primitive.ObjectID, price float64)
}

type pricesRepository struct{}

func NewPricesRepository() PricesRepository{
	return &pricesRepository{}
}

func (pr * pricesRepository) GetPrice(ctx context.Context, date int64, groupId primitive.ObjectID) (*domain.Prices, error) {

	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(PricesCollectionName)

	var result domain.Prices
	filter := bson.D{{"date", date}, {"groupid", groupId}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
		return &result, err
	}
	return &result, nil
}

//adds or changing price on the date and group
func (pr * pricesRepository) AddPrice(ctx context.Context, date int64, groupId primitive.ObjectID, price float64) {

	database, err := InitDb(ctx)
	if err != nil {
		return // nil, err
	}

	collection := database.Collection(PricesCollectionName)

	result, err := pr.GetPrice(ctx, date, groupId)
	if err != nil {
		fmt.Println(err)
	}

	//price on date not found
	//need to add new price
	if result == nil {
		result = domain.NewPrices()
		result.ID = primitive.NewObjectID()
		result.Date = date
		result.GroupID = groupId
		result.Price = price

		insertResult, err := collection.InsertOne(ctx, result)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(insertResult)

		//price on date founded
		//need to change existing price
	} else {
		if result.Price != price {
			searchFilter := bson.M{"_id": result.ID}
			updateFilter := bson.M{"$set": bson.M{"price": price}}
			updateResult, err := collection.UpdateOne(ctx, searchFilter, updateFilter)
			if err != nil {
				fmt.Printf("update fail %v\n", err)
			}
			fmt.Println(updateResult)
		}
	}
	//fmt.Println(collection)
	err = database.Client().Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
		return //nil, err
	}
}

func (pr * pricesRepository) ChangePrice(ctx context.Context, date int64, groupId primitive.ObjectID, price float64){
	database, err := InitDb(ctx)
	if err != nil {
		return // nil, err
	}

	collection := database.Collection(PricesCollectionName)

	result, err := pr.GetPrice(ctx, date, groupId)
	if err != nil {
		fmt.Println(err)
	}

	if result != nil {
		if result.Price != price {
			searchFilter := bson.M{"_id": result.ID}
			updateFilter := bson.M{"$set": bson.M{"price": price}}
			updateResult, err := collection.UpdateOne(ctx, searchFilter, updateFilter)
			if err != nil {
				fmt.Printf("update fail %v\n", err)
			}
			fmt.Println(updateResult)
		}
	}
	err = database.Client().Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
		return //nil, err
	}
}
