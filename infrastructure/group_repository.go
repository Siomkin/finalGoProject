package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main/domain"
)

//TODO create interface and struct

type GroupRepository interface {
	GetGroupByName(ctx context.Context, groupName string) (*domain.Group, error)
	AddGroup(ctx context.Context, groupName string) error
}

type groupRepository struct {}

func NewGroupRepository() GroupRepository{
	return &groupRepository{}
}

func (gr *groupRepository) GetGroupByName(ctx context.Context, groupName string) (*domain.Group, error) {
	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(GroupCollectionName)
	fmt.Println(collection)
	var result domain.Group

	filter := bson.D{{"name", groupName}}
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		fmt.Println("returning nil & error")
		return nil, err
	}

	return &result, nil
}

func (gr *groupRepository) AddGroup(ctx context.Context, groupName string) error{
	database, err := InitDb(ctx)
	if err != nil {
		return err
	}

	collection := database.Collection(GroupCollectionName)

	var result *domain.Group

	result, err = gr.GetGroupByName(ctx, groupName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if result == nil {

		result = domain.NewGroup()
		result.ID = primitive.NewObjectID()
		result.Name = groupName

		insertResult, err := collection.InsertOne(ctx, result)
		if err != nil {
			return err
		}
		fmt.Println(insertResult)

	}
	return nil
}
