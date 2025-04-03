package database

import (
	"context"
	"time"

	"github.com/Vantuan1606/app-buff/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	Client *mongo.Client
}

var mongoService *MongoService

func NewMongoService() *MongoService {
	if mongoService != nil && mongoService.Client != nil {

		return mongoService
	}

	config := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.Uri))

	if err != nil {

		logrus.WithError(err).Error("Mongo connect fail")
		panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logrus.WithError(err).Error("Ping mongo fail")
		panic(err)
	}

	mongoService = &MongoService{
		Client: client,
	}

	return mongoService
}
