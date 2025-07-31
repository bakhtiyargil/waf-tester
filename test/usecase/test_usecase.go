package usecase

import (
	"context"
	"time"
	"waf-tester/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testUseCase struct {
	testRepo       domain.TestRepository
	contextTimeout time.Duration
}

func NewCatUseCase(u domain.TestRepository, to time.Duration) domain.TestUseCase {
	return &testUseCase{
		testRepo:       u,
		contextTimeout: to,
	}
}

func (cat *testUseCase) InsertOne(c context.Context, m *domain.Test) (*domain.Test, error) {
	ctx, cancel := context.WithTimeout(c, cat.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now().In(time.FixedZone("AZT", 4*60*60))

	res, err := cat.testRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (cat *testUseCase) FindOne(c context.Context, id string) (*domain.Test, error) {
	ctx, cancel := context.WithTimeout(c, cat.contextTimeout)
	defer cancel()

	res, err := cat.testRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (cat *testUseCase) DeleteOne(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, cat.contextTimeout)
	defer cancel()

	err := cat.testRepo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
