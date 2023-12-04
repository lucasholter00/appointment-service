package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func Connect() {
	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	Database = c.Database("AppointmentService")
	if err != nil {
		panic(err)
	}
	fmt.Println("App is connected to MongoDB")
}
