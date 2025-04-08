package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Vantuan1606/app-test/user"

)

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
	Password string             `json:"password" bson:"password,omiempty"`
	
}

type UsersRequest struct {
	ID string `json:"id"`
}

type IUserUsecase interface {
	List(context.Context, *user.ListUserInput) ([]*User, error)
	GetUser(context.Context, string) (*User, error)
}

type IUserRepo interface {
	FindOne(context.Context, map[string]interface{}) (*User, error)
	Find(context.Context, map[string]interface{}, ...*options.FindOptions) ([]*User, error)
}
