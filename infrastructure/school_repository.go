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
type SchoolRepository interface{

	GetSchools(ctx context.Context) ([] *domain.School, error)
	AddSchool(ctx context.Context, name string)(primitive.ObjectID, error)
	GetSchoolByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error)
	GetSchoolByName(ctx context.Context, name string) (*domain.School, error)
}

type schoolRepository struct {}

func NewSchoolRepository() SchoolRepository {
	return &schoolRepository{}
}


//Python
//def get_schools(self):
//self.cur.execute(r"SELECT * FROM mydb.school")
//rows = self.cur.fetchall()
//return rows
//
func (sr *schoolRepository) GetSchools(ctx context.Context) ([]*domain.School, error){
	database, err := InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	collection := database.Collection(SchoolsCollectionName)

	var schools []*domain.School

	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for cur.Next(ctx){
		var elem domain.School
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		schools = append(schools, &elem)
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {

		fmt.Println(err)
		return schools, err
	}
	return schools, nil
}

func (sr *schoolRepository) GetSchoolByID(ctx context.Context, id primitive.ObjectID) (*domain.School, error){
	database, err := InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	collection := database.Collection(SchoolsCollectionName)
	filter := bson.D{{"_id", id}}
	var school domain.School

	err = collection.FindOne(ctx, filter).Decode(school)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
		return &school, err
	}

	return &school, nil
}

func (sr *schoolRepository) AddSchool(ctx context.Context, name string) (*domain.School, error){

	//var emptyVal primitive.ObjectID
	school, err := sr.GetSchoolByName(ctx, name)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if school == nil {

		database, err := InitDb(ctx)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		collection := database.Collection(SchoolsCollectionName)

		newschool := domain.NewSchool()
		newschool.ID = primitive.NewObjectID()
		newschool.Name = name

		res, err := collection.InsertOne(ctx, newschool)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println(res)
		return newschool, nil
	}
	return nil, err
}

func (sr *schoolRepository) GetSchoolByName(ctx context.Context, name string) (*domain.School, error){
	database, err := InitDb(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	collection := database.Collection(SchoolsCollectionName)
	filter := bson.D{{"name", name}}
	var school domain.School
	err = collection.FindOne(ctx, filter).Decode(school)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = database.Client().Disconnect(ctx)
	if err != nil {
		fmt.Println(err)
		return &school, err
	}

	return &school, nil
}
