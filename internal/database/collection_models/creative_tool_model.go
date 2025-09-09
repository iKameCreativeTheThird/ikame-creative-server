package collectionmodels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreativeTool struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Team     string             `bson:"team"`
	ToolName string             `bson:"tool_name"`
	Type     string             `bson:"type"`
	Point    []int              `bson:"point"`
}
