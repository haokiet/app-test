package repo

import (
	"context"

	"github.com/Vantuan1606/app-test/config"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/service/database"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUserRepo struct {
	c *mongo.Collection
}

func NewMongoUserRepo() domain.IUserRepo {
	c := config.GetConfig()
	mongo := database.NewMongoService()

	return &mongoUserRepo{
		c: mongo.Client.Database(c.Mongo.Database).Collection("User"),
	}
}

func (m *mongoUserRepo) Find(ctx context.Context, conditions map[string]interface{}, options ...*options.FindOptions) ([]*domain.User, error) {
	cursor, err := m.c.Find(ctx, conditions, options...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var user *domain.User
		if err := cursor.Decode(&user); err != nil {
			logrus.WithError(err).Error("[DECODE] fail at ", user.ID)
			return nil, err
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *mongoUserRepo) FindOne(ctx context.Context, conditions map[string]interface{}) (*domain.User, error) {
	var user *domain.User
	if err := m.c.FindOne(ctx, conditions).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
