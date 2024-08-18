package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"goland01/db"
	"goland01/model"
	"goland01/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConsumeKafka() {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Panic(err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}
	defer partitionConsumer.Close()

	client, ctx, cancel, err := db.ConnectMongo()
	if err != nil {
		log.Panic(err)
	}
	defer cancel()
	defer client.Disconnect(ctx)

	for message := range partitionConsumer.Messages() {
		handleMessage(message.Value, client, ctx)
	}
}

func handleMessage(value []byte, client *mongo.Client, ctx context.Context) {
	var collectionName string
	var event interface{}

	switch model.Service {
	case "person":
		event = &model.Person{}
		collectionName = "person"
	case "document":
		event = &model.Person{}
		collectionName = "document"
	default:
		log.Printf("unknown collection type: %s", model.Service)
		return
	}

	rawMessage := string(value)
	cleanedMessage := utils.CleanJsonString(rawMessage)

	err := json.Unmarshal([]byte(cleanedMessage), event)
	if err != nil {
		log.Printf("error deserializing message: %v, %v", err, cleanedMessage)
		return
	}

	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection(collectionName)

	// Obtain the id from the event
	var id int
	switch v := event.(type) {
	case *model.Person:
		id = v.ID
	default:
		log.Printf("unexpected type: %T", v)
		return
	}

	filter := bson.M{"id": id}

	var existingDoc bson.M
	err = collection.FindOne(ctx, filter).Decode(&existingDoc)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("error checking if document exists: %v", err)
		return
	}

	if err == mongo.ErrNoDocuments {
		// Document doesn't exist, insert a new one
		_, err = collection.InsertOne(ctx, event)
		if err != nil {
			log.Printf("error inserting to MongoDB: %v", err)
		} else {
			fmt.Println("Document inserted to MongoDB:", event)
		}
	} else {
		// Document exists, update it
		update := bson.M{"$set": event}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Printf("error updating MongoDB document: %v", err)
		} else {
			fmt.Println("Document updated in MongoDB:", event)
		}
	}
}