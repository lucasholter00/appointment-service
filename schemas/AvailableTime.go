package schemas

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AvailableTime struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Dentist_id primitive.ObjectID `bson:"Dentist_id"`
	Start_time primitive.DateTime `bson:"Start_time"`
	End_time   primitive.DateTime `bson:"End_time"`
}
