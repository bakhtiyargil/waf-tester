package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Test struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Host           string             `bson:"host" json:"host" validate:"required"`
	Path           string             `bson:"path" json:"path" validate:"required"`
	Method         string             `bson:"method" json:"method" validate:"required"`
	ResponseStatus int                `bson:"response_status" json:"response_status" validate:"required"`
	ResponseBdy    string             `bson:"response_bdy" json:"response_bdy" validate:"required"`
	ResponseTime   float64            `bson:"response_time" json:"response_time" validate:"required"`
	TestID         string             `bson:"test_id" json:"test_id" validate:"required"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

type TestRepository interface {
	FindOne(ctx context.Context, id string) (*Test, error)
	InsertOne(ctx context.Context, u *Test) (*Test, error)
	DeleteOne(ctx context.Context, id string) error
}

type TestUseCase interface {
	FindOne(ctx context.Context, id string) (*Test, error)
	InsertOne(ctx context.Context, u *Test) (*Test, error)
	DeleteOne(ctx context.Context, id string) error
}
