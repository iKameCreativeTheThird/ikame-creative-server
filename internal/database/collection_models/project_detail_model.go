package collectionmodels

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectDetail struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID int                `bson:"id"`
	Project   string             `bson:"project"`
	Research  string             `bson:"research"`
	Art       string             `bson:"art"`
	Concept   string             `bson:"concept"`
	Video     string             `bson:"video"`
	Pla       string             `bson:"pla"`
	UA        string             `bson:"ua"`
}

func InstertNewProjectDetailToDatabase(client *mongo.Client, url, dbName, collName string, projectDetail *ProjectDetail) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)
	_, err := collection.InsertOne(ctx, projectDetail)
	return err
}
