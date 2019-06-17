package infrastructure

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDb(ctx context.Context) (*mongo.Database, error){
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