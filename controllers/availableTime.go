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
			panic(err)
		}

		go CreateAvailableTime(payload, client)
	})

	if tokenCreate.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenGet := client.Subscribe("grp20/availabletimes/get", byte(0), func(c mqtt.Client, m mqtt.Message) {
		var payload schemas.AvailableTime
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		go GetAllAvailableTimesWithDentistID(payload.Dentist_id, client)
	})

	if tokenGet.Error() != nil {
		panic(tokenCreate.Error())
	}

	tokenDelete := client.Subscribe("grp20/dentist/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.AvailableTime
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		go DeleteAvailableTime(payload.ID, client)
	})

	if tokenDelete.Error() != nil {
		panic(tokenCreate.Error())
	}

}

// fungerar
func CreateAvailableTime(payload schemas.AvailableTime, client mqtt.Client) bool {

	col := getAvailableTimesCollection()
	// Hash the password using Bcrypt
	fmt.Printf("Start")

	result, err := col.InsertOne(context.TODO(), payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Registered availableTime with dentistID: %v \n", result.InsertedID)
	return true

}

//fungerar men hur ska vi förmedla den

// getAllInstancesWithDentistID retrieves all documents in a collection with a matching dentist_id
func GetAllAvailableTimesWithDentistID(dentistID string, client mqtt.Client) (*mongo.Cursor, error) {

	col := getAvailableTimesCollection()
	filter := bson.D{{Key: "dentist_id", Value: dentistID}}

	collection, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for i := 0; i < collection.RemainingBatchLength(); i++ {
		fmt.Printf("hej")
	}
	return collection, nil
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

	fmt.Printf("Deleted Time id: %v \n", ID)
	return true

}

func getAvailableTimesCollection() *mongo.Collection {
	col := database.Database.Collection("AvailableTimes")
	return col
}
