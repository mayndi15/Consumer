package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func ConnectMongo() (*mongo.Client, context.Context, context.CancelFunc, error) {
    clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    client, err := mongo.Connect(ctx, clientOptions)
    return client, ctx, cancel, err
}
