package collectionmodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WeeklyOrder struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	DateFrom time.Time          `bson:"date_from"`
	DateTo   time.Time          `bson:"date_to"`
	Goal     string             `bson:"goal"`
	Strategy string             `bson:"strategy"`
	Project  string             `bson:"project"`
	CPP      int                `bson:"cpp"`
	Icon     int                `bson:"icon"`
	Banner   int                `bson:"banner"`
	PLA      int                `bson:"pla"`
	Video    int                `bson:"video"`
}
