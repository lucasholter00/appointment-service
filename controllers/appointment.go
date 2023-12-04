package controllers

import (
	"Group20/appointment-service/schemas"
	"Group20/appointment-service/database"
	"context"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitialiseAppointment (client mqtt.Client) {

//CREATE
    tokenCreate := client.Subscribe("grp20/req/appointment/delete", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Appointment
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        go CancelAppointment(payload._id, client)
    })
    if tokenCreate.Error() != nil {
        panic(tokenCreate.Error())
    }
}

func CancelAppointment(id, client mqtt.Client) bool{
    var message string
    var code string
    var returnVal bool

    col := getAppointmentCollection()
    user := &schemas.Appointment
}


func getAppointmentCollection() *mongo.Collection {
    col := database.Database.Collection("Appointments")
    return col
}
