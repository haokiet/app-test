package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"

	"github.com/Vantuan1606/app-test/domain"
	"github.com/Vantuan1606/app-test/hashtag"
)

type hashtagUsecase struct {
	hashtagRepo    domain.IHashtagRepo
	contextTimeout time.Duration
}

func NewHashtagUsecase(pr domain.IHashtagRepo, timeOut time.Duration) domain.IHashtagUsecase {
	return &hashtagUsecase{
		hashtagRepo:    pr,
		contextTimeout: timeOut,
	}
}

func (hgs *hashtagUsecase) GetHashtag(c context.Context, hashtagID string) (*domain.Hashtag, error) {
	ctx, cancel := context.WithTimeout(c, hgs.contextTimeout)
	defer cancel()

	hashtagObjectID, err := primitive.ObjectIDFromHex(hashtagID)
	log.Println("hashtagObjectID", hashtagObjectID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("Invalid ID")
	}

	conditions := bson.M{"_id": hashtagObjectID}
	hashtag, err := hgs.hashtagRepo.FindOne(ctx, conditions)

	if err != nil {
		logrus.WithField("Data", fmt.Sprintf("%v", conditions)).WithError(err).Error("hashtag not found")
		return nil, err
	}

	return hashtag, nil
}

func (hgs *hashtagUsecase) List(c context.Context, input *hashtag.ListHashtagInput) ([]*domain.Hashtag, error) {
	ctx, cancel := context.WithTimeout(c, hgs.contextTimeout)
	defer cancel()

	conditions := bson.M{}

	options := options.Find()
	options.SetLimit(int64(*input.Limit)) // Giới hạn số lượng user trả về

	hashtags, err := hgs.hashtagRepo.Find(ctx, conditions, options)
	if err != nil {
		logrus.WithError(err).Error("Get list hashtag failed")
		return nil, err
	}

	return hashtags, nil
}
