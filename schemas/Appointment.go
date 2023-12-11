package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

type Appointment struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"` 
    Patient_id primitive.ObjectID `bson:"patient_id"` 
    Dentist_id primitive.ObjectID `bson:"dentist_id"`
    Start_time primitive.DateTime `bson:"start_time"`
    End_time   primitive.DateTime `bson:"end_time"`
    Clinic_id  primitive.ObjectID `bson:"Clinic_id"`
}
