package infrastructure

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Connection interface{
	InitDb(ctx context.Context) (*mongo.Database, error)
	CloseDb(ctx context.Context, db *mongo.Database) error
}

type conn struct{}
func NewConnection() Connection{
	return &conn{}
}


func (cn *conn) InitDb(ctx context.Context) (*mongo.Database, error){
	// Create client
	client, err := mongo.NewClient(options.Client().ApplyURI(Credentials.Host + ":" + Credentials.Port))
	if err != nil {
		return nil, err
	}
	// Create connect
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(Credentials.DBName) ,nil
}

func (cn *conn) CloseDb(ctx context.Context, db *mongo.Database) error {
	err := db.Client().Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return err
}