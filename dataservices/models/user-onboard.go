package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserOnboard struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string
	Phone   string
	Gender  string
	SteamID string
}
