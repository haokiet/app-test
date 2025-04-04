package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/user"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type userUsecase struct {
	userRepo       domain.IUserRepo
	contextTimeout time.Duration
}

func NewUserUsecase(pr domain.IUserRepo, timeOut time.Duration) domain.IUserUsecase {
	return &userUsecase{
		userRepo:       pr,
		contextTimeout: timeOut,
	}
}

func (us *userUsecase) GetUser(c context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	log.Println("userID", userObjectID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Invalid ID")
	}

	conditions := bson.M{"_id": userObjectID}
	user, err := us.userRepo.FindOne(ctx, conditions)

	if err != nil {
		logrus.WithField("Data", fmt.Sprintf("%v", conditions)).WithError(err).Error("User not found")
		return nil, err
	}

	return user, nil
}

// func (us *userUsecase) List(c context.Context, input *user.ListUserInput) ([]*domain.User, error) {
// 	ctx, cancel := context.WithTimeout(c, us.contextTimeout)
// 	defer cancel()

// 	conditions := bson.M{}

// 	options := options.Find()

// 	options.SetSkip(int64(*input.Offset))
// 	options.SetLimit(int64(*input.Limit))

// 	users, err := us.userRepo.Find(ctx, conditions, options)

// 	if err != nil {
// 		logrus.WithError(err).Error("Get list user failed")
// 		return nil, err
// 	}

// 	return users, nil

// }

func (us *userUsecase) List(c context.Context, input *user.ListUserInput) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, us.contextTimeout)
	defer cancel()

	conditions := bson.M{}

	options := options.Find()
	options.SetLimit(int64(*input.Limit)) // Giới hạn số lượng user trả về

	users, err := us.userRepo.Find(ctx, conditions, options)
	if err != nil {
		logrus.WithError(err).Error("Get list user failed")
		return nil, err
	}

	return users, nil
}
