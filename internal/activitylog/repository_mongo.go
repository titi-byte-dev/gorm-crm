package activitylog

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName     = "activity_logs"
	logRetentionDays   = 90
	logRetentionSecs   = logRetentionDays * 24 * 60 * 60
	indexTimeout       = 10 * time.Second
	queryTimeout       = 5 * time.Second
	defaultResultLimit = 50
)

var _ Repository = (*mongoRepository)(nil)

type mongoRepository struct {
	col *mongo.Collection
}

// NewMongoRepository cria o repositório e garante os índices necessários.
// Índices no MongoDB são criados explicitamente — ao contrário do PostgreSQL,
// não há FK automáticos que criem índices.
func NewMongoRepository(db *mongo.Database) Repository {
	col := db.Collection(collectionName)

	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "entity_type", Value: 1}, {Key: "entity_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(logRetentionSecs),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), indexTimeout)
	defer cancel()
	col.Indexes().CreateMany(ctx, indexes) //nolint:errcheck — índices criados se não existirem

	return &mongoRepository{col: col}
}

func (r *mongoRepository) Save(log *Log) error {
	if log.ID.IsZero() {
		log.ID = primitive.NewObjectID()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	_, err := r.col.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("save activity log: %w", err)
	}
	return nil
}

func (r *mongoRepository) FindByEntity(entityType EntityType, entityID string, limit int) ([]*Log, error) {
	if limit <= 0 {
		limit = defaultResultLimit
	}
	filter := bson.M{"entity_type": entityType, "entity_id": entityID}
	// Sort por _id descrescente = mais recente primeiro (ObjectID tem timestamp embutido)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(int64(limit))
	return r.find(filter, opts)
}

func (r *mongoRepository) FindByUser(userID string, limit int) ([]*Log, error) {
	if limit <= 0 {
		limit = defaultResultLimit
	}
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(int64(limit))
	return r.find(filter, opts)
}

func (r *mongoRepository) find(filter bson.M, opts *options.FindOptions) ([]*Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find activity logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*Log
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("decode activity logs: %w", err)
	}
	return logs, nil
}
