package repository

import (
	"context"
	"manager/internal/model"
)

type RequestsRepositoryInterface interface {
	Save(ctx context.Context, record *model.CrackHashRecord) (*model.CrackHashRecord, error)
	GetByID(ctx context.Context, ID string) (*model.CrackHashRecord, error)
	Update(ctx context.Context, record *model.CrackHashRecord) (*model.CrackHashRecord, error)
}

type RequestsMongoRepository struct {
	// mongo connection
	// logger
}
