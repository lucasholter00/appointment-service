package controllers

import (
	"Group20/appointment-service/database"
	"Group20/appointment-service/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeAvailableTimes(client mqtt.Client) {

	tokenCreate := client.Subscribe("grp20/dentist/post", byte(0), func(c mqtt.Client, m mqtt.Message) {
		var payload schemas.AvailableTime

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			var message string
			message = "{\"message\": \"Bad request\",\"code\": \"400\"}"
			client.Publish("grp20/res/availabletime/create", 0, false, message)
		} else {
			go CreateAvailableTime(payload, client, false)
		}
	})

	if tokenCreate.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenGet := client.Subscribe("grp20/req/timeSlots/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
		var payload schemas.AvailableTime
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			message := "{\"Message\": \"Bad request\",\"Code\": \"400\"}"
			client.Publish("grp20/res/timeslots/get", 0, false, message)
		} else {
			go GetAllAvailableTimesWithDentistID(payload.Dentist_id, client)
		}
	})

	if tokenGet.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenDelete := client.Subscribe("grp20/req/dentist/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.AvailableTime
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			message := "{\"Message\": \"Bad request\",\"Code\": \"400\"}"
			client.Publish("grp20/res/dentist/delete", 0, false, message)
		} else {
			if !exist(payload) {
				message := "{\"Message\": \"Not found!\",\"Code\": \"404\"}"
				client.Publish("grp20/res/dentist/delete", 0, false, message)
			} else {
				go DeleteAvailableTime(payload.ID, client)
			}
		}
	})

	if tokenDelete.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenInternalMigrate := client.Subscribe("appointmentservice/internal/migrate", byte(0), func(c mqtt.Client, m mqtt.Message) {
		var payload schemas.AvailableTime

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			fmt.Printf("malformed payload!")
		} else {
			go CreateAvailableTime(payload, client, true)
		}
	})

	if tokenInternalMigrate.Error() != nil {
		panic(tokenInternalMigrate.Error())
	}

}

func CreateAvailableTime(payload schemas.AvailableTime, client mqtt.Client, internal bool) bool {
	var message string

	if exist(payload) {
		var message string
		message = "{\"Message\": \"An identical available time already exist!\",\"Code\": \"409\"}"
		client.Publish("grp20/res/availabletime/create", 0, false, message)
		return false
	}

	if payload.Start_time > payload.End_time {
		message = "{\"Message\": \"End time must be after the start time\",\"Code\": \"409\"}"
		client.Publish("grp20/res/availabletime/create", 0, false, message)
		return false
	}

	col := getAvailableTimesCollection()

	result, err := col.InsertOne(context.TODO(), payload)

	if internal == false {
		if err != nil {
			log.Fatal(err)
			message = "{\"Message\": \"An error occurred\",\"Code\": \"500\"}"
			client.Publish("grp20/res/availabletime/create", 0, false, message)
			return false
		}

		message = "{\"Message\": \"Available time created\",\"Code\": \"201\"}"
		fmt.Printf("Registered availableTime with dentistID: %v \n", result.InsertedID)
		client.Publish("grp20/res/availabletime/create", 0, false, message)
		return true
	} else {
		if err != nil {
			return false
		} else {
			client.Publish("internal/res", 0, false, "data migrated!")
			return true
		}
	}
}

// getAllInstancesWithDentistID retrieves all documents in a collection with a matching dentist_id
func GetAllAvailableTimesWithDentistID(dentistID primitive.ObjectID, client mqtt.Client) bool {

	col := getAvailableTimesCollection()
	filter := bson.D{{Key: "dentist_id", Value: dentistID}}

	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		message := "{\"Message\": \"An error occurred\",\"Code\": \"500\"}"
		client.Publish("grp20/res/timeslots/get", 0, false, message)
		return false
	}

	defer cursor.Close(context.TODO())

	var availableTimes []schemas.AvailableTime

	for cursor.Next(context.TODO()) {
		var availableTime schemas.AvailableTime
		if err := cursor.Decode(&availableTime); err != nil {
			message := "{\"Message\": \"An error occurred while decoding results\",\"Code\": \"500\"}"
			client.Publish("grp20/res/timeslots/get", 0, false, message)
			return false
		}
		availableTimes = append(availableTimes, availableTime)
	}

	if err := cursor.Err(); err != nil {
		message := "{\"Message\": \"An error occurred\",\"Code\": \"500\"}"
		client.Publish("grp20/res/timeslots/get", 0, false, message)
		return false
	}

	// Convert the responseMap to JSON
	resultJSON, err := json.Marshal(availableTimes)

	if err != nil {
		message := "{\"Message\": \"An error occurred while converting to JSON\",\"Code\": \"500\"}"
		client.Publish("grp20/res/timeslots/get", 0, false, message)
		return false
	}

	fmt.Printf(string(resultJSON))
	client.Publish("grp20/res/timeslots/get", 0, false, string(resultJSON))

	return true
}

// deletes an availableTime entirely, will be performed by dentists
func DeleteAvailableTime(ID primitive.ObjectID, client mqtt.Client) bool {

	col := getAvailableTimesCollection()
	filter := bson.M{"_id": ID}
	result, err := col.DeleteOne(context.TODO(), filter)
	_ = result

	if err != nil {
		log.Fatal(err)
		return false
	}

	msg := schemas.AvailableTime{
		ID: ID,
	}

	if result.DeletedCount == 0 {
		document, err := json.Marshal(msg)

		if err != nil {
			message := "{\"Message\": \"Internal server error!\",\"Code\": \"500\"}"
			client.Publish("grp20/res/dentist/delete", 0, false, message)
			return false
		}

		client.Publish("appointmentservice/internal/delete", 0, false, document)

		return false
	} else {
		fmt.Printf("Deleted Time id: %v \n", ID)
		message := "{\"Message\": \"Available time deleted!\",\"Code\": \"200\"}"
		client.Publish("grp20/res/dentist/delete", 0, false, message)

		return true

	}
}

func exist(payload schemas.AvailableTime) bool {
	col := getAvailableTimesCollection()

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

func getAvailableTimesCollection() *mongo.Collection {
	col := database.Database.Collection("AvailableTimes")
	return col
}
