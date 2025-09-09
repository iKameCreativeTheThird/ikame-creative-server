package collectionmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
