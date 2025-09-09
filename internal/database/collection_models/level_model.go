package collectionmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Level struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Team       string             `bson:"team"`
	LevelPoint []int              `bson:"levelPoint"`
}
