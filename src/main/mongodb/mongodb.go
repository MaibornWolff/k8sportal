package mongodb

import (
	"context"
	"fmt"
	"k8sportal/model"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout = 10 * time.Second
)

var mongodbdatabase = "k8sportal"
var mongodbcollection = "portal-services"

func Connect(ctx context.Context, mongodbHost string) (*mongo.Client, error) {
	timedContext, cancelTimedContext := context.WithTimeout(ctx, connectionTimeout)
	defer cancelTimedContext()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbHost))
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %w", err)
	}

	err = client.Connect(timedContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize client: %w", err)
	}

	err = client.Ping(timedContext, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to ping MongoDB: %w", err)
	}

	log.Info().Msg("Connected to MongoDB")

	return client, nil
}

func GetAllServices(mongoClient *mongo.Client) ([]*model.Service, error) {
	var services []*model.Service

	ctx := context.Background() //TODO get context from function call

	db := mongoClient.Database(mongodbdatabase)
	collection := db.Collection(mongodbcollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &services)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return services, nil
}
