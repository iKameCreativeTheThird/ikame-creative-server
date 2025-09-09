package collectionmodels

import (
	"time"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompletedTask struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	TaskID     string             `bson:"id"`
	TaskName   string             `bson:"task_name"`
	AssigneeID string             `bson:"assignee_id"`
	Tool       []string           `bson:"tool"`
	Level      int                `bson:"level"`
	TaskType   string             `bson:"task_type"`
	Project    string             `bson:"project"`
	Team       string             `bson:"team"`
	DoneDate   time.Time          `bson:"done_date"`
}

//	{
//	  done_date:
//	  {
//	    $gte: new Date("2025-09-01T00:00:00.000Z"),
//	    $lte: new Date("2025-09-02T23:59:59.999Z")
//	  }
//	}
//
// GetCompletedTasksByDateRange retrieves completed tasks between two dates (inclusive)
func GetCompletedTasksByDateRange(ctx context.Context, collection *mongo.Collection, startDate, endDate time.Time) ([]CompletedTask, error) {
	filter := bson.M{
		"done_date": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}
	opts := options.Find()
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []CompletedTask
	for cursor.Next(ctx) {
		var task CompletedTask
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		results = append(results, task)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
