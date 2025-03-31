package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddOrderType(w http.ResponseWriter, r *http.Request) {
	var orderType models.OrderType
	if err := json.NewDecoder(r.Body).Decode(&orderType); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.OrderTypesCollection)
	orderType.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.Background(), orderType)
	if err != nil {
		http.Error(w, "Error adding order type", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderType)
}

func GetOrderTypes(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.OrderTypesCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching order types", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var orderTypes []models.OrderType
	if err = cursor.All(context.Background(), &orderTypes); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderTypes)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	// Get order type name
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err := orderTypeCollection.FindOne(context.Background(), bson.M{"_id": order.OrderTypeID}).Decode(&orderType)
	if err != nil {
		http.Error(w, "Order type not found", http.StatusBadRequest)
		return
	}

	order.OrderType = orderType.Name

	collection := db.GetCollection(db.OrderCollection)
	order.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), order)
	if err != nil {
		http.Error(w, "Error adding order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.OrderCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
