package controllers

import (
	"Group20/appointment-service/database"
	"Group20/appointment-service/schemas"
	"context"
	"encoding/json"
	"fmt"

	//"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitialiseAppointment(client mqtt.Client) {

	//CREATE
	tokenCancel := client.Subscribe("grp20/req/appointment/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Appointment
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			//send mqtt message 400 bad request
			fmt.Printf("400 - Bad request")
		} else {
			go CancelAppointment(payload.ID, client)
		}
	})
	if tokenCancel.Error() != nil {
		panic(tokenCancel.Error())
	}

	tokenCreate := client.Subscribe("grp20/req/appointment/create", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Appointment
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			//send mqtt message 400 bad request
			fmt.Printf("400 - Bad request")
		} else {
			go CreateAppointment(payload, client)
		}
	})
	if tokenCreate.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenGetAllForUser := client.Subscribe("grp20/req/appointment/get", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Appointment
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			//send mqtt message 400 bad request
			fmt.Printf("400 - Bad request")
		} else {
			go GetAllForUser(payload.Patient_id, client)
		}
	})
	if tokenGetAllForUser.Error() != nil {
		panic(tokenGetAllForUser.Error())
	}

	tokenDelete := client.Subscribe("appointmentservice/internal/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Appointment
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			//send mqtt message 400 bad request or ignore due to internal?
			fmt.Printf("400 - Bad request")
		} else {
			go DeleteAppointment(payload.ID, client)
		}
	})
	if tokenDelete.Error() != nil {
		panic(tokenDelete.Error())
	}

}

func DeleteAppointment(id primitive.ObjectID, client mqtt.Client) bool {
	var message string
	var code string
	var returnVal bool

	col := getAppointmentCollection()
	filter := bson.M{"_id": id}
	result, err := col.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	if result.DeletedCount == 1 {

		message = `{"message": "User deleted"`
		code = "200"

		message = AddCodeStringJson(message, code)
	} else {
		message = `{"message": "User not found"`
		code = "404"
		message = AddCodeStringJson(message, code)
	}

	client.Publish("grp20/res/dentist/delete", byte(0), false, message)
	returnVal = true
	return returnVal
}

func GetAllForUser(id primitive.ObjectID, client mqtt.Client) bool {

	var message string
	var code string
	var returnVal bool

	col := getAppointmentCollection()
	filter := bson.M{"patient_id": id}

	cursor, err := col.Find(context.TODO(), filter)

	if err != nil {
		message = `{"message": "Error"`
		code = "500"
	}

	defer cursor.Close(context.TODO())

	var appointments []schemas.Appointment

	for cursor.Next(context.TODO()) {
		var appointment schemas.Appointment

		if err := cursor.Decode(&appointment); err != nil {
			message = `{"message": "Error"`
			code = "500"
			panic(err)
		}

		appointments = append(appointments, appointment)
	}
	jsonData, err := json.Marshal(appointments)
	fmt.Printf(string(jsonData))
	returnVal = true

	code = "200"
	message = AddCodeStringJson(string(jsonData), code)

	client.Publish("grp20/res/appointment/get", 0, false, message)

	return returnVal
}

func CreateAppointment(payload schemas.Appointment, client mqtt.Client) bool {
	var message string
	var code string
	var returnVal bool

	if exist(payload) {
		//send conflict - http status code
		return false
	}

	col := getAppointmentCollection()

	result, err := col.InsertOne(context.TODO(), payload)
	_ = result

	if err != nil {
		code = "500"
		message = `{"message": "Appointmend could not be created"`
		client.Publish("grp20/res/appointment/create", 0, false, message)
		panic(err)
	}

	message = `{"message": "Appointment booked"`
	code = "200"
	returnVal = true

	message = AddCodeStringJson(message, code)

	client.Publish("grp20/res/appointment/create", 0, false, message)
	return returnVal
}

func CancelAppointment(id primitive.ObjectID, client mqtt.Client) bool {
	var message string
	var code string
	var returnVal bool
	appointment := &schemas.Appointment{}

	col := getAppointmentCollection()
	filter := bson.M{"_id": id}
	data := col.FindOne(context.TODO(), filter)

	if data.Err() == mongo.ErrNoDocuments {
		//send 404 message
		return false
	}

	result, err := col.DeleteOne(context.TODO(), filter)

	if err != nil {
		panic(err)
	}

	if result.DeletedCount == 1 {

		code = "200"

		data.Decode(appointment)

		availableTime := schemas.AvailableTime{
			Dentist_id: appointment.ID,
			Start_time: appointment.Start_time,
			End_time:   appointment.End_time,
		}

		jsonData, err := json.Marshal(availableTime)

		if err != nil {
			panic(err)
		}

		message = string(jsonData)

		fmt.Printf("Här")
		client.Publish("appointmentservice/internal/migrate", 0, false, message)

		message = `{"message": "Appointment Canceled"`
		message = AddCodeStringJson(message, code)

		returnVal = true

	} else {

		code = "404"
		message = `{"message": "Appointment not found"`
		message = AddCodeStringJson(message, code)

		returnVal = false
	}

	client.Publish("grp20/res/appointment/delete", 0, false, message)
	return returnVal
}

func exist(payload schemas.Appointment) bool {
	col := getAppointmentCollection()

	filter := bson.M{
		"Dentist_id": payload.Dentist_id,
		"Start_time": payload.Start_time,
		"End_time":   payload.End_time,
	}

	count, err := col.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}

	return count > 0
}

func getAppointmentCollection() *mongo.Collection {
	col := database.Database.Collection("Appointments")
	return col
}
