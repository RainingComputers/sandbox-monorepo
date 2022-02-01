package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}

func (LHS Post) IsEqual(RHS Post) bool {
	LHS.ID = RHS.ID
	return LHS == RHS
}
