package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"main/domain"
)


type ChildrenRepository interface{
	GetChildrenById(ctx context.Context, childID string) (*domain.Children, error)
	GetChildrenByUserID(ctx context.Context, userID string) ([] *domain.Children, error)
	GetChildrenByNameAndUserID(ctx context.Context, childName string, userID string) (*domain.Children, error)
	AddChild(ctx context.Context, userID string, groupID string, childName string) (*domain.Children, error)
	DeleteChild(ctx context.Context, userID string, childName string) error
}

type childrenRepository struct{}


func NewChildrenRepository() ChildrenRepository{
	return &childrenRepository{}
}

func (cr *childrenRepository) GetChildrenByNameAndUserID(ctx context.Context, childName string, userID string) (*domain.Children, error) {
	cn := NewConnection()
	database, err := cn.InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cn.CloseDb(ctx, database)

	collection := database.Collection(ChildrenCollectionName)

	var result domain.Children
	_id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"name", childName}, {"userid", _id}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//def viewchildren(self, user_id):
//self.cur.execute(r"SELECT * FROM mydb.children WHERE userid = %s", (user_id,))
//rows = self.cur.fetchall()
//return rows
func (cr *childrenRepository) GetChildrenById(ctx context.Context, childID string) (*domain.Children, error){
	cn := NewConnection()
	database, err := cn.InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cn.CloseDb(ctx, database)

	collection := database.Collection(ChildrenCollectionName)

	var result domain.Children
	filter := bson.D{{"_id", childID}}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


func (cr *childrenRepository) GetChildrenByUserID(ctx context.Context, userID string) ([] *domain.Children, error) {
	cn := NewConnection()
	database, err := cn.InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cn.CloseDb(ctx, database)

	collection := database.Collection(ChildrenCollectionName)

	var children [] *domain.Children
	_id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"userid", _id}}

	result, err := collection.Find(ctx, filter)
	if err != nil {
		if err != mongo.ErrNoDocuments{
			fmt.Println(err)
			return nil, err
		} else {
			return nil, nil
		}
	}

	for result.Next(ctx) {
		var elem domain.Children
		err := result.Decode(&elem)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		children = append(children, &elem)
	}

	return children, nil
}


//def add_child(self, child_name, user_id, group_id):
//self.cur.execute(r"INSERT INTO mydb.children VALUES (NUll, %s, %s, %s)", (child_name, group_id,  user_id,))
//self.conn.commit()
func (cr *childrenRepository) AddChild(ctx context.Context, userID string, groupID string, childName string) (*domain.Children, error){
	//var emptyVal primitive.ObjectID
	cn := NewConnection()
	database, err := cn.InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cn.CloseDb(ctx, database)

	collection := database.Collection(ChildrenCollectionName)

	child, err := cr.GetChildrenByNameAndUserID(ctx, childName, userID)

	if err != nil{
		if err != mongo.ErrNoDocuments{
			fmt.Println(err)
			return nil, err
		}
	}
	if child == nil {
		child := domain.NewChildren()
		child.ID = primitive.NewObjectID()
		child.UserID, _ = primitive.ObjectIDFromHex(userID)
		child.GroupID, _ = primitive.ObjectIDFromHex(groupID)
		child.Name = childName

		insertResult, err := collection.InsertOne(ctx, child)
		if err != nil {
			return nil, err
		}
		fmt.Println(insertResult)
	}

	return child, nil
}

//def delete_child(self, name, userid):
//self.cur.execute(r"DELETE FROM mydb.children WHERE name=%s and userid=%s", (name, userid,))
//self.conn.commit()
func (cr *childrenRepository) DeleteChild(ctx context.Context, userID string, childName string) error{

	return nil
}