package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	var existingUser models.User
	err = collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь с таким email уже существует"})
		return
	}

	err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь с таким username уже существует"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя", http.StatusInternalServerError)
		return
	}

	tokenString, err := utils.GenerateToken(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка при получении токена", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	var user models.User
	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	var updateData struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
		FullName    string `json:"fullName"`
		Email       string `json:"email"`
		BirthDate   string `json:"birthDate"`
		Residence   string `json:"residence"`
		Education   string `json:"education"`
		BirthPlace  string `json:"birthPlace"`
		HomeAddress string `json:"homeAddress"`
	}

	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateData.OldPassword))
	if err != nil {
		http.Error(w, "Неверный старый пароль", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"fullName":    updateData.FullName,
			"email":       updateData.Email,
			"birthDate":   updateData.BirthDate,
			"residence":   updateData.Residence,
			"education":   updateData.Education,
			"birthPlace":  updateData.BirthPlace,
			"homeAddress": updateData.HomeAddress,
			"password":    string(hashedPassword),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": userID}, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении профиля", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Профиль успешно обновлен"})
}
