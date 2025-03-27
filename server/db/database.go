package db

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName = "diplome"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
)

// GetMongoClient возвращает экземпляр клиента MongoDB (синглтон)
func GetMongoClient() *mongo.Client {
	clientOnce.Do(func() {
		// В реальном проекте URI лучше брать из конфига
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		
		var err error
		clientInstance, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Проверяем подключение
		err = clientInstance.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}
		
		log.Println("Successfully connected to MongoDB!")
	})
	return clientInstance
}

// GetCollection возвращает коллекцию из базы данных
func GetCollection(collectionName string) *mongo.Collection {
	return GetMongoClient().Database(DatabaseName).Collection(collectionName)
}

// CloseConnection закрывает соединение с MongoDB
func CloseConnection() {
	if clientInstance != nil {
		if err := clientInstance.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed")
	}
}