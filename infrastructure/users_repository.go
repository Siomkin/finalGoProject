package infrastructure

import (
	"context"
	"crypto/md5"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"main/domain"
)

//func GetUserByID(id string) *domain.User{
//
//
//}
//

type UsersRepository interface{
	UserLogin(ctx context.Context, login string, pass string) bool
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CreateUser(ctx context.Context, login string, pass string) (*domain.User, error)
}

type usersRepository struct{}

func NewUsersRepository() UsersRepository{
	return &usersRepository{}
}

func (u *usersRepository) UserLogin(ctx context.Context, login string, pass string) bool{
	us, err := u.GetUserByLogin(ctx, login)
	if err != nil {
		log.Fatal(err)
	}

	if us != nil{
		h := md5.New()

		loginPass := string(h.Sum([]byte(pass)))
		if loginPass == us.Pass{
			return true
		}
	}
	return false
}

func (u *usersRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error){
	var result domain.User

	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	collection := database.Collection(UsersCollectionName)
	filter := bson.D{{"name", login}}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connection to MongoDB closed.")

	return &result, nil
}

func (u *usersRepository) CreateUser(ctx context.Context, login string, pass string) (*domain.User, error){

	database, err := InitDb(ctx)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	collection := database.Collection(UsersCollectionName)

	userForInsert := domain.NewUser()
	userForInsert.SetNewName(login)
	userForInsert.Pass = string(h.Sum([]byte(pass)))
	userForInsert.ID = primitive.NewObjectID()

	insertResult, err := collection.InsertOne(ctx, userForInsert)
	if err != nil {
		return nil, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	err = database.Client().Disconnect(ctx)
	if err != nil {
		return nil, err
	}

	return userForInsert, nil
}