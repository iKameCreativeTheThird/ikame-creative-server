package collectionmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeeklyTarget struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Team     string             `bson:"team"`
	Point    int                `bson:"point"`
	DateFrom time.Time          `bson:"date_from"`
	DateTo   time.Time          `bson:"date_to"`
}
