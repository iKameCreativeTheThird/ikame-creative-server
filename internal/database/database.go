package db_handler

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURI := os.Getenv("MONGO_URI")
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	// Optionally, ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	log.Println("Connected to MongoDB!")
	return nil
}

type PerformancePoint struct {
	TotalPerformancePoint     float64 `bson:"total_performance_point"`
	TotalCreativeProcessPoint float64 `bson:"total_creative_process_point"`
	TotalCreativeTaskPoint    float64 `bson:"total_creative_task_point"`
	TotalBasePoint            float64 `bson:"total_base_point"`
}

func GetPerformancePointForIndividual(uri, dbName, collName, memberEmail string, startDate, endDate time.Time) (*PerformancePoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database(dbName).Collection(collName)

	// Pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "assignee", Value: memberEmail},
			{Key: "done_date", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "creative-tool"},
			{Key: "localField", Value: "tool"},
			{Key: "foreignField", Value: "tool_name"},
			{Key: "as", Value: "tool_info"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "performance_point", Value: "$level"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "tool_factor", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$eq", Value: bson.A{bson.D{{Key: "$size", Value: "$tool_info"}}, 0}},
					}},
					{Key: "then", Value: 0},
					{Key: "else", Value: 1}, // placeholder to keep example short
				}},
			}},
			{Key: "creative_process_point", Value: 0}, // placeholder
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "creative_task_point", Value: bson.D{
				{Key: "$multiply", Value: bson.A{"$performance_point", "$tool_factor"}},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "base_point", Value: bson.D{
				{Key: "$add", Value: bson.A{
					bson.D{{Key: "$subtract", Value: bson.A{"$performance_point", "$creative_task_point"}}},
					"$creative_process_point",
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_performance_point", Value: bson.D{{Key: "$sum", Value: "$performance_point"}}},
			{Key: "total_creative_process_point", Value: bson.D{{Key: "$sum", Value: "$creative_process_point"}}},
			{Key: "total_creative_task_point", Value: bson.D{{Key: "$sum", Value: "$creative_task_point"}}},
			{Key: "total_base_point", Value: bson.D{{Key: "$sum", Value: "$base_point"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []PerformancePoint
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}

func GetPerformancePointForTeam(uri, dbName, collName, team string, startDate, endDate time.Time) (*PerformancePoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database(dbName).Collection(collName)

	// Pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "team", Value: team},
			{Key: "done_date", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "creative-tool"},
			{Key: "localField", Value: "tool"},
			{Key: "foreignField", Value: "tool_name"},
			{Key: "as", Value: "tool_info"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "performance_point", Value: "$level"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "tool_factor", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$eq", Value: bson.A{bson.D{{Key: "$size", Value: "$tool_info"}}, 0}},
					}},
					{Key: "then", Value: 0},
					{Key: "else", Value: 1}, // placeholder to keep example short
				}},
			}},
			{Key: "creative_process_point", Value: 0}, // placeholder
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "creative_task_point", Value: bson.D{
				{Key: "$multiply", Value: bson.A{"$performance_point", "$tool_factor"}},
			}},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "base_point", Value: bson.D{
				{Key: "$add", Value: bson.A{
					bson.D{{Key: "$subtract", Value: bson.A{"$performance_point", "$creative_task_point"}}},
					"$creative_process_point",
				}},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total_performance_point", Value: bson.D{{Key: "$sum", Value: "$performance_point"}}},
			{Key: "total_creative_process_point", Value: bson.D{{Key: "$sum", Value: "$creative_process_point"}}},
			{Key: "total_creative_task_point", Value: bson.D{{Key: "$sum", Value: "$creative_task_point"}}},
			{Key: "total_base_point", Value: bson.D{{Key: "$sum", Value: "$base_point"}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []PerformancePoint
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return &results[0], nil
}
