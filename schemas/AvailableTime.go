package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type AvailableTime struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Dentist_id string             `bson:"dentist_id"`
	Start_time primitive.DateTime `bson:"startTime"`
	End_time   primitive.DateTime `bson:"endTime"`
}
