package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.UsersCollection)

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.GenerateToken(user.ID.Hex())
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

func CheckToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		fmt.Println("Неверный токен:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
