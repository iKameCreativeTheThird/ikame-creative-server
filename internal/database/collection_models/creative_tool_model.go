package collectionmodels

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreativeTool struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Team     string             `bson:"team"`
	ToolName string             `bson:"tool_name"`
	Type     string             `bson:"type"`
	Point    []int              `bson:"point"`
}

func GetCreativeToolByTeam(client *mongo.Client, dbName, collectionName, team string) (*CreativeTool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collectionName)
	var tool CreativeTool
	err := collection.FindOne(ctx, bson.M{"team": team}).Decode(&tool)
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

func AddCreativeTool(client *mongo.Client, dbName, collectionName string, tool *CreativeTool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(ctx, tool)
	return err
}

func UpdateCreativeTool(client *mongo.Client, dbName, collectionName string, tool *CreativeTool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.UpdateOne(ctx,
		bson.M{"team": tool.Team, "tool_name": tool.ToolName},
		bson.M{"$set": bson.M{"type": tool.Type, "point": tool.Point}},
	)
	return err
}

func DeleteCreativeTool(client *mongo.Client, dbName, collectionName, team, toolName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"team": team, "tool_name": toolName})
	return err
}
