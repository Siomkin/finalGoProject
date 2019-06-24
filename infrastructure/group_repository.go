package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"main/domain"
)

//TODO create interface and struct

type GroupRepository interface {
	GetGroupByName(ctx context.Context, groupName string) (*domain.Group, error)
	AddGroup(ctx context.Context, groupName string, SchoolID string)  (*domain.Group, error)
	GetGroups(ctx context.Context) ([] *domain.Group, error)
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

func (gr *groupRepository) AddGroup(ctx context.Context, groupName string, SchoolID string)  (*domain.Group, error){
	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(GroupCollectionName)

	var result *domain.Group

	result, err = gr.GetGroupByName(ctx, groupName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if result == nil {

		result = domain.NewGroup()
		result.ID = primitive.NewObjectID()
		result.Name = groupName
		result.SchoolID, _ = primitive.ObjectIDFromHex(SchoolID)

		insertResult, err := collection.InsertOne(ctx, result)
		if err != nil {
			return nil, err
		}
		fmt.Println(insertResult)

	}
	return result, nil
}

func (gr *groupRepository) GetGroups(ctx context.Context) ([] *domain.Group, error) {
	database, err := InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	collection := database.Collection(GroupCollectionName)

	var groups []*domain.Group

	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for cur.Next(ctx){
		var elem domain.Group
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, &elem)
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {

		fmt.Println(err)
		return groups, err
	}
	return groups, nil
}
