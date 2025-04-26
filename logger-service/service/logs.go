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
	ID bson.ObjectID `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (ls *LoggerService) Insert(entry LogEntry) error {
	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		ID: bson.NewObjectID(),
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

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error getting all logs:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var entry LogEntry
		if err := cursor.Decode(&entry); err != nil {
			log.Println("Error decoding log entry: ", err)
			continue
		}
		logs = append(logs, &entry)
	}

	return  logs, nil
}

func (ls *LoggerService) GetOne(id string) (*LogEntry, error) {
	//Convert string to ObjectID
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting string to ObjectID:", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&entry)
	if err != nil {
		log.Println("Error getting single log entry:", err)
		return nil, err
	}

	return &entry, nil
}

func (ls *LoggerService) UpdateOne(entryToUpdate LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	result, err := collection.UpdateOne(ctx, bson.M{"_id": entryToUpdate.ID}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: entryToUpdate.Name},
			{Key: "data", Value: entryToUpdate.Data},
			{Key: "updated_at", Value: time.Now()},
			{Key: "created_at", Value: entryToUpdate.CreatedAt},
		}},
	})
	if err != nil {
		log.Printf("Error updating log entry %s: %v", entryToUpdate.ID, err)
		return nil, err
	}

	return result, nil
}

func (ls *LoggerService) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := ls.mClient.Database(ls.dbName).Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection logs:", err)
		return err
	}

	return nil
}