package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Vantuan1606/app-test/hashtag"
)

type Hashtag struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name,omitempty"`
	Status int                `json:"status" bson:"status,omiempty"`
}

type HashtagsRequest struct {
	ID string `json:"id"`
}

type IHashtagUsecase interface {
	List(context.Context, *hashtag.ListHashtagInput) ([]*Hashtag, error)
	GetHashtag(context.Context, string) (*Hashtag, error)
}

type IHashtagRepo interface {
	FindOne(context.Context, map[string]interface{}) (*Hashtag, error)
	Find(context.Context, map[string]interface{}, ...*options.FindOptions) ([]*Hashtag, error)
}
