package collectionmodels

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Member struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	MemberID string             `bson:"id"`
	Name     string             `bson:"name"`
	YOB      int                `bson:"yob"`
	Email    string             `bson:"email"`
	Role     string             `bson:"role"`
	Team     string             `bson:"team"`
}

func UpdateMemberToDataBase(client *mongo.Client, url, dbName, collName string, member *Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"id": member.MemberID},
		bson.M{"$set": bson.M{
			"name":  member.Name,
			"yob":   member.YOB,
			"email": member.Email,
			"role":  member.Role,
			"team":  member.Team,
		}},
	)
	return err
}

func InsertMemberToDataBase(client *mongo.Client, url, dbName, collName string, member *Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)

	_, err := collection.InsertOne(ctx, member)
	return err
}

func DeleteMemberInDataBase(client *mongo.Client, url, dbName, collName, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)
	_, err := collection.DeleteOne(ctx, bson.M{"id": memberID})
	return err
}

func GetMemberByEmail(client *mongo.Client, url, dbName, collName, email string) (*Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)
	var member Member
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func GetAllMembers(client *mongo.Client, url, dbName, collName string) ([]*Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(collName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var members []*Member
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}
	return members, nil
}
