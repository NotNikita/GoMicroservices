package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LoggerService struct {
	mClient *mongo.Client
	dbName string
}

func NewLoggerService(client *mongo.Client) (*LoggerService, error) {
	return &LoggerService{
		mClient: client,
		dbName: "udemy_logs",
	}, nil
}

// Log object that would be stored in the MongoDB
type LogEntry struct {
	ID string `bson:"_id", omitempty json:"id", omitempty`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"created_at" json:"created_at"`
}

func (ls *LoggerService) Insert(entry LogEntry) error {
	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	_, err := collection.InsertOne(nil, LogEntry{
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry:", err)
		return err
	}

	return nil
}

func (ls *LoggerService) GetAll(limit_opt ...int64) ([]*LogEntry, error) {
	var limit int64 = 30
	if len(limit_opt) > 0 {
		limit = limit_opt[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	opts := options.Find().SetSort(bson.D{bson.E{Key: "created_at",Value: -1}}).SetLimit(limit)

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error getting all logs:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		if err := cursor.Decode(entry); err != nil {
			log.Println("Error decoding log entry: ", err)
			return nil, err
		}
		logs = append(logs, &entry)
	}

	return  logs, nil
}