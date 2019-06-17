package infrastructure

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"main/domain"
)

//TODO create interface and struct


type ChildrenRepository interface{
	GetChildrenById(ctx context.Context, childID primitive.ObjectID) (*domain.Children, error)
	GetChildrenByUserID(ctx context.Context, userID primitive.ObjectID) ([] *domain.Children, error)
	AddChild(ctx context.Context, userID primitive.ObjectID, groupID primitive.ObjectID, childName string) (primitive.ObjectID, error)
	DeleteChild(ctx context.Context, userID primitive.ObjectID, childName string) error
}

type childrenRepository struct{}

func NewchildrenRepository() ChildrenRepository{
	return &childrenRepository{}
}

//def viewchildren(self, user_id):
//self.cur.execute(r"SELECT * FROM mydb.children WHERE userid = %s", (user_id,))
//rows = self.cur.fetchall()
//return rows
func (cr *childrenRepository) GetChildrenById(ctx context.Context, childID primitive.ObjectID) (*domain.Children, error){
	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(ChildrenCollectionName)

	var result domain.Children
	filter := bson.D{{"_id", childID}}
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


func (cr *childrenRepository) GetChildrenByUserID(ctx context.Context, userID primitive.ObjectID) ([] *domain.Children, error) {

	return nil, nil
}


//def add_child(self, child_name, user_id, group_id):
//self.cur.execute(r"INSERT INTO mydb.children VALUES (NUll, %s, %s, %s)", (child_name, group_id,  user_id,))
//self.conn.commit()
func (cr *childrenRepository) AddChild(ctx context.Context, userID primitive.ObjectID, groupID primitive.ObjectID, childName string) (primitive.ObjectID, error){

	return nil
}

//def delete_child(self, name, userid):
//self.cur.execute(r"DELETE FROM mydb.children WHERE name=%s and userid=%s", (name, userid,))
//self.conn.commit()
func (cr *childrenRepository) DeleteChild(ctx context.Context, userID primitive.ObjectID, childName string) error{

	return nil
}