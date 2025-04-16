package repository

import (
	"context"
	"errors"
	"manager/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type RequestsMongoRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewRequestsMongoRepository(db *mongo.Database, logger *zap.Logger) *RequestsMongoRepository {
	return &RequestsMongoRepository{
		collection: db.Collection("crack_hash_requests"),
		logger:     logger,
	}
}

func (r *RequestsMongoRepository) Save(ctx context.Context, record *model.CrackHashRecord) (*model.CrackHashRecord, error) {
	record.CreatedAt = time.Now().UTC()
	record.Answers = nil
	record.ProgressInPercents = 0.0

	_, err := r.collection.InsertOne(ctx, record)
	if err != nil {
		r.logger.Error("Failed to insert CrackHashRecord", zap.Error(err))
		return nil, err
	}

	r.logger.Info("CrackHashRecord saved", zap.String("id", record.ID))
	return record, nil
}

func (r *RequestsMongoRepository) GetByID(ctx context.Context, ID string) (*model.CrackHashRecord, error) {
	var result model.CrackHashRecord

	err := r.collection.FindOne(ctx, bson.M{"id": ID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("CrackHashRecord not found", zap.String("id", ID))
			return nil, nil
		}
		r.logger.Error("Failed to find CrackHashRecord", zap.String("id", ID), zap.Error(err))
		return nil, err
	}

	return &result, nil
}

func (r *RequestsMongoRepository) Update(ctx context.Context, record *model.CrackHashRecord) (*model.CrackHashRecord, error) {
	filter := bson.M{"id": record.ID}
	update := bson.M{"$set": record}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.CrackHashRecord

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		r.logger.Error("Failed to update CrackHashRecord", zap.String("id", record.ID), zap.Error(err))
		return nil, err
	}

	return &updated, nil
}
