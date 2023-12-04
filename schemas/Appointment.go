package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Appointment struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"` 
    patient_id primitive.ObjectID `bson:"patient_id"` 
    dentist_id primitive.ObjectID `bson:"dentist_id"`
    start_time primitive.DateTime `bson:"start_time"`
    end_time   primitive.DateTime `bson:"end_time"`
}
