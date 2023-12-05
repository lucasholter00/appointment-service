package schemas

import (
	"fmt"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AvailableTime struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Dentist_id primitive.ObjectID `bson:"Dentist_id" validate:"required,hexadecimal,len=24"`
	Start_time primitive.DateTime `bson:"Start_time" validate:"required"`
	End_time   primitive.DateTime `bson:"End_time" validate:"required,gtefield=StartTime"`
}

var validate = validator.New()

func (at *AvailableTime) Validate() error {
	if err := validate.Struct(at); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}
	return nil
}
