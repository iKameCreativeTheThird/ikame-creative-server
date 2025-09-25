package db_handler

import (
	"context"
	"log"
	"os"
	collectionmodels "performance-dashboard-backend/internal/database/collection_models"
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
	StartWeek                 time.Time `bson:"_id"`
	TotalPerformancePoint     float64   `bson:"total_performance_point"`
	TotalCreativeProcessPoint float64   `bson:"total_creative_process_point"`
	TotalCreativeTaskPoint    float64   `bson:"total_creative_task_point"`
	TotalBasePoint            float64   `bson:"total_base_point"`
}

func GetPerformancePoint(uri, dbName, collName, indentifier string, startDate, endDate time.Time, isTeam bool) ([]*PerformancePoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// defer client.Disconnect(ctx)
	collection := client.Database(dbName).Collection(collName)

	var identifierKey string
	if isTeam {
		identifierKey = "team"
	} else {
		identifierKey = "assignee_id"
	}

	// Pipeline
	// Build aggregation pipeline
	pipeline := mongo.Pipeline{
		// 1. Filter by date range
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: identifierKey, Value: indentifier},
				{Key: "done_date", Value: bson.D{
					{Key: "$gte", Value: startDate},
					{Key: "$lte", Value: endDate},
				}},
			}},
		},

		// 2. Add week_start (truncate to Monday)
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "week_start", Value: bson.D{
					{Key: "$dateTrunc", Value: bson.D{
						{Key: "date", Value: "$done_date"},
						{Key: "unit", Value: "week"},
						{Key: "binSize", Value: 1},
						{Key: "timezone", Value: "UTC"},
						{Key: "startOfWeek", Value: "monday"},
					}},
				}},
			}},
		},

		// 3. Lookup creative-tool
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "creative-tool"},
				{Key: "localField", Value: "tool"},
				{Key: "foreignField", Value: "tool_name"},
				{Key: "as", Value: "tool_info"},
			}},
		},

		// 4. Separate tool type "t" and "q"
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "performance_point", Value: "$level"},
				{Key: "tool_points_t", Value: bson.D{
					{Key: "$map", Value: bson.D{
						{Key: "input", Value: bson.D{
							{Key: "$filter", Value: bson.D{
								{Key: "input", Value: "$tool_info"},
								{Key: "as", Value: "t"},
								{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$t.type", "t"}}}},
							}},
						}},
						{Key: "as", Value: "t"},
						{Key: "in", Value: bson.D{
							{Key: "tool_name", Value: "$$t.tool_name"},
							{Key: "point", Value: bson.D{
								{Key: "$let", Value: bson.D{
									{Key: "vars", Value: bson.D{
										{Key: "idx", Value: bson.D{{Key: "$subtract", Value: bson.A{"$level", 1}}}},
										{Key: "arr", Value: "$$t.point"},
									}},
									{Key: "in", Value: bson.D{
										{Key: "$cond", Value: bson.A{
											bson.D{{Key: "$lt", Value: bson.A{"$$idx", bson.D{{Key: "$size", Value: "$$arr"}}}}},
											bson.D{{Key: "$arrayElemAt", Value: bson.A{"$$arr", "$$idx"}}},
											bson.D{{Key: "$arrayElemAt", Value: bson.A{"$$arr", 0}}},
										}},
									}},
								}},
							}},
						}},
					}},
				}},
				{Key: "tool_points_q", Value: bson.D{
					{Key: "$map", Value: bson.D{
						{Key: "input", Value: bson.D{
							{Key: "$filter", Value: bson.D{
								{Key: "input", Value: "$tool_info"},
								{Key: "as", Value: "t"},
								{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$t.type", "q"}}}},
							}},
						}},
						{Key: "as", Value: "t"},
						{Key: "in", Value: bson.D{
							{Key: "tool_name", Value: "$$t.tool_name"},
							{Key: "point", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$$t.point", 0}}}},
						}},
					}},
				}},
			}},
		},

		// 5. Calculate tool_factor and creative_process_point
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "tool_factor", Value: bson.D{
					{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{bson.D{{Key: "$size", Value: "$tool_points_t"}}, 0}}}},
						{Key: "then", Value: 0},
						{Key: "else", Value: bson.D{
							{Key: "$reduce", Value: bson.D{
								{Key: "input", Value: bson.D{
									{Key: "$map", Value: bson.D{
										{Key: "input", Value: "$tool_points_t"},
										{Key: "as", Value: "tp"},
										{Key: "in", Value: bson.D{
											{Key: "$cond", Value: bson.A{
												bson.D{{Key: "$lt", Value: bson.A{"$$tp.point", 1}}},
												"$$tp.point",
												1,
											}},
										}},
									}},
								}},
								{Key: "initialValue", Value: 1},
								{Key: "in", Value: bson.D{{Key: "$multiply", Value: bson.A{"$$value", "$$this"}}}},
							}},
						}},
					}},
				}},
				{Key: "creative_process_point", Value: bson.D{
					{Key: "$sum", Value: bson.D{
						{Key: "$map", Value: bson.D{
							{Key: "input", Value: "$tool_points_q"},
							{Key: "as", Value: "tp"},
							{Key: "in", Value: "$$tp.point"},
						}},
					}},
				}},
			}},
		},

		// 6. Calculate creative_task_point
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "creative_task_point", Value: bson.D{{Key: "$multiply", Value: bson.A{"$performance_point", "$tool_factor"}}}},
			}},
		},

		// 7. Calculate base_point
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "base_point", Value: bson.D{
					{Key: "$add", Value: bson.A{
						bson.D{{Key: "$subtract", Value: bson.A{"$performance_point", "$creative_task_point"}}},
						"$creative_process_point",
					}},
				}},
			}},
		},
		// 8. Group by week
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$week_start"},
				{Key: "total_performance_point", Value: bson.D{{Key: "$sum", Value: "$performance_point"}}},
				{Key: "total_creative_process_point", Value: bson.D{{Key: "$sum", Value: "$creative_process_point"}}},
				{Key: "total_creative_task_point", Value: bson.D{{Key: "$sum", Value: "$creative_task_point"}}},
				{Key: "total_base_point", Value: bson.D{{Key: "$sum", Value: "$base_point"}}},
			}},
		},

		// 9. Sort by week
		bson.D{
			{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*PerformancePoint
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}

func GetMembersByTeam(uri, dbName, collName string, team string) ([]*collectionmodels.Member, error) {
	// Example body request
	// 	{
	//     "teams": ["Art Creative"]
	// 	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// defer client.Disconnect(ctx)
	collection := client.Database(dbName).Collection(collName)

	// if team == "" or null return all collection
	if team == "" {
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		var results []*collectionmodels.Member
		if err = cursor.All(ctx, &results); err != nil {
			return nil, err
		}
		return results, nil
	}

	filter := bson.M{
		"team": team,
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*collectionmodels.Member
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
