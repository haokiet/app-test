package repo

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Vantuan1606/app-test/config"
	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/service/database"
)

type mongoHashtagRepo struct {
	c *mongo.Collection
}

func NewMongoHashtagRepo() domain.IHashtagRepo {
	c := config.GetConfig()
	mongo := database.NewMongoService()

	return &mongoHashtagRepo{
		c: mongo.Client.Database(c.Mongo.Database).Collection("Hashtag"),
	}
}

func (m *mongoHashtagRepo) Find(ctx context.Context, conditions map[string]interface{}, options ...*options.FindOptions) ([]*domain.Hashtag, error) {
	cursor, err := m.c.Find(ctx, conditions, options...)
	if err != nil {
		logrus.WithError(err).Error("[FIND] Query failed")
		return nil, err
	}
	defer cursor.Close(ctx)

	var hashtags []*domain.Hashtag
	for cursor.Next(ctx) {
		var hashtag *domain.Hashtag
		if err := cursor.Decode(&hashtag); err != nil {
			logrus.WithError(err).Error("[DECODE] Failed to decode hashtag")
			return nil, err
		}
		hashtags = append(hashtags, hashtag)
	}

	if err := cursor.Err(); err != nil {
		logrus.WithError(err).Error("[FIND] Cursor encountered an error")
		return nil, err
	}

	return hashtags, nil
}

func (m *mongoHashtagRepo) FindOne(ctx context.Context, conditions map[string]interface{}) (*domain.Hashtag, error) {

	var hashtag *domain.Hashtag
	if err := m.c.FindOne(ctx, conditions).Decode(&hashtag); err != nil {
		logrus.WithError(err).Error("[FINDONE] Query failed or decode error")
		return nil, err
	}
	return hashtag, nil
}
