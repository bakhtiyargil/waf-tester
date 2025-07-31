package mongo

import (
	"context"
	"waf-tester/domain"
	"waf-tester/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

const (
	collectionName = "test"
)

func NewMongoRepository(DB mongo.Database) domain.TestRepository {
	return &mongoRepository{DB, DB.Collection(collectionName)}
}

func (m *mongoRepository) InsertOne(ctx context.Context, test *domain.Test) (*domain.Test, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, test)
	if err != nil {
		return test, err
	}

	return test, nil
}

func (m *mongoRepository) FindOne(ctx context.Context, id string) (*domain.Test, error) {
	var (
		cat domain.Test
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &cat, err
	}

	err = m.Collection.FindOne(ctx, bson.M{"_id": idHex}).Decode(&cat)
	if err != nil {
		return &cat, err
	}

	return &cat, nil
}

func (m *mongoRepository) DeleteOne(ctx context.Context, id string) error {
	var (
		err error
	)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": idHex}
	_, err = m.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
