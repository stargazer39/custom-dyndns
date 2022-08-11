package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DDNSEntry struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Hostname string             `json:"hostname" bson:"hostname"`
	IP       string             `json:"ip" bson:"ip"`
	Added    time.Time          `json:"added" bson:"added"`
	Updated  time.Time          `json:"updated" bson:"updated"`
}
