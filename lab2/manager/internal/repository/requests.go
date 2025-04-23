package repository

import (
	"context"
	"errors"
	"fmt"
	"manager/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type RequestsMongoRepository struct {
	requestsCollection *mongo.Collection
	subTasksCollection *mongo.Collection
	logger             *zap.Logger
}

func NewRequestsMongoRepository(db *mongo.Database, logger *zap.Logger) *RequestsMongoRepository {
	return &RequestsMongoRepository{
		logger:             logger,
		requestsCollection: db.Collection("crack_hash_requests"),
		subTasksCollection: db.Collection("subtasks"),
	}
}

func (r *RequestsMongoRepository) Save(ctx context.Context, record *model.CrackHashRecord) (*model.CrackHashRecord, error) {
	record.CreatedAt = time.Now().UTC()
	record.Answers = nil
	record.ProgressInPercents = 0.0

	_, err := r.requestsCollection.InsertOne(ctx, record)
	if err != nil {
		r.logger.Error("Failed to insert CrackHashRecord", zap.Error(err))
		return nil, err
	}

	r.logger.Info("CrackHashRecord saved", zap.String("id", record.ID))
	return record, nil
}

func (r *RequestsMongoRepository) SaveSubTask(ctx context.Context, subTask *model.SubTask) (*model.SubTask, error) {
	subTask.CreatedAt = time.Now().UTC()
	subTask.Answers = nil
	subTask.Progress = 0.0

	_, err := r.subTasksCollection.InsertOne(ctx, subTask)
	if err != nil {
		r.logger.Error("Failed to insert CrackHashSubTask", zap.Error(err))
		return nil, err
	}

	r.logger.Info("CrackHashSubTask saved", zap.String("id", subTask.ID))
	return subTask, nil
}

func (r *RequestsMongoRepository) MarkSubTaskAsFinished(ctx context.Context, subTaskID string, answers []string) error {
	filter := bson.M{"id": subTaskID}
	update := bson.M{
		"$set": bson.M{
			"status":   "FINISHED",
			"progress": 1.0,
			"answers":  answers,
		},
	}
	_, err := r.subTasksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to mark SubTask as finished", zap.String("id", subTaskID), zap.Error(err))
	}
	return err
}

func (r *RequestsMongoRepository) AreAllSubTasksFinished(ctx context.Context, taskID string) (bool, error) {
	filter := bson.M{
		"task_id": taskID,
		"status":  bson.M{"$ne": "FINISHED"},
	}

	count, err := r.subTasksCollection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count unfinished subtasks", zap.String("task_id", taskID), zap.Error(err))
		return false, err
	}

	return count == 0, nil
}

func (r *RequestsMongoRepository) MarkMainTaskAsFinished(ctx context.Context, taskID string, answers []string) error {
	now := time.Now().UTC()
	filter := bson.M{"id": taskID}
	update := bson.M{
		"$set": bson.M{
			"status":      "READY",
			"answers":     answers,
			"finished_at": now,
		},
	}

	_, err := r.requestsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to mark main task as finished", zap.String("id", taskID), zap.Error(err))
	}

	return err
}

func (r *RequestsMongoRepository) CollectAllAnswersByTaskID(ctx context.Context, taskID string) ([]string, error) {
	filter := bson.M{"taskid": taskID}
	cursor, err := r.subTasksCollection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to find subtasks for collecting answers", zap.String("task_id", taskID), zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.SubTask
	if err := cursor.All(ctx, &results); err != nil {
		r.logger.Error("Failed to decode subtask documents", zap.Error(err))
		return nil, err
	}
	fmt.Println("RESULTS: ", results)
	var allAnswers []string
	for _, task := range results {
		if task.Answers != nil {
			allAnswers = append(allAnswers, task.Answers...)
		}
	}

	return allAnswers, nil
}

func (r *RequestsMongoRepository) GetByID(ctx context.Context, ID string) (*model.CrackHashRecord, error) {
	var result model.CrackHashRecord

	err := r.requestsCollection.FindOne(ctx, bson.M{"id": ID}).Decode(&result)
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

	err := r.requestsCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		r.logger.Error("Failed to update CrackHashRecord", zap.String("id", record.ID), zap.Error(err))
		return nil, err
	}

	return &updated, nil
}

func (r *RequestsMongoRepository) UpdateSubTaskProgress(ctx context.Context, subTaskID string, progress float64) error {
	filter := bson.M{"id": subTaskID}
	update := bson.M{
		"$set": bson.M{
			"progress": progress,
		},
	}

	_, err := r.subTasksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update SubTask progress", zap.String("id", subTaskID), zap.Error(err))
	}
	return err
}

func (r *RequestsMongoRepository) GetSubTaskByID(ctx context.Context, subTaskID string) (*model.SubTask, error) {
	var subTask model.SubTask

	err := r.subTasksCollection.FindOne(ctx, bson.M{"id": subTaskID}).Decode(&subTask)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("SubTask not found", zap.String("id", subTaskID))
			return nil, nil
		}
		r.logger.Error("Failed to find SubTask", zap.String("id", subTaskID), zap.Error(err))
		return nil, err
	}

	return &subTask, nil
}
