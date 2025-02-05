package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
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
	var updateData struct {
		FullName          string `json:"fullName,omitempty"`
		LastName          string `json:"lastName,omitempty"`
		FirstName         string `json:"firstName,omitempty"`
		MiddleName        string `json:"middleName,omitempty"`
		Education         string `json:"education,omitempty"`
		Phone         string `json:"phone,omitempty"`
		BirthDate         string `json:"birthDate,omitempty"`
		BirthPlace        string `json:"birthPlace,omitempty"`
		HomeAddress       string `json:"homeAddress,omitempty"`
		WorkPlace         string `json:"workPlace,omitempty"`
		JobTitle          string `json:"jobTitle,omitempty"`
		OldPassword       string `json:"oldPassword,omitempty"`
		NewPassword       string `json:"newPassword,omitempty"`
		PassportData      string `json:"passportData,omitempty"`
		SNILS             string `json:"snils,omitempty"`
		AgreeToProcessing bool   `json:"agreeToProcessing,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из токена
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}
	
	userID := claims.UserID

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	// Получаем текущие данные пользователя
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(userID)
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Создаем объект для обновления
	update := bson.M{}

	if updateData.FullName != "" {
		update["fullName"] = updateData.FullName
	}
	if updateData.LastName != "" {
		update["lastname"] = updateData.LastName
	}
	if updateData.FirstName != "" {
		update["firstname"] = updateData.FirstName
	}
	if updateData.MiddleName != "" {
		update["middlename"] = updateData.MiddleName
	}
	if updateData.Education != "" {
		update["education"] = updateData.Education
	}
	if updateData.Phone != "" {
		update["phone"] = updateData.Phone
	}
	if updateData.BirthDate != "" {
		update["birthdate"] = updateData.BirthDate
	}
	if updateData.BirthPlace != "" {
		update["birthplace"] = updateData.BirthPlace
	}
	if updateData.HomeAddress != "" {
		update["homeaddress"] = updateData.HomeAddress
	}
	if updateData.WorkPlace != "" {
		update["workplace"] = updateData.WorkPlace
	}
	if updateData.JobTitle != "" {
		update["jobtitle"] = updateData.JobTitle
	}
	if updateData.PassportData != "" {
		update["passportdata"] = updateData.PassportData
	}
	if updateData.SNILS != "" {
		update["snils"] = updateData.SNILS
	}
	if updateData.AgreeToProcessing {
		update["agreetoprocessing"] = updateData.AgreeToProcessing
	}

	// Если пользователь хочет изменить пароль
	if updateData.OldPassword != "" && updateData.NewPassword != "" {
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

		update["password"] = string(hashedPassword)
	}

	// Обновляем только те поля, которые были переданы
	if len(update) > 0 {
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": update})
		if err != nil {
			http.Error(w, "Ошибка при обновлении данных", http.StatusInternalServerError)
			return
		}
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Данные успешно обновлены"})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Ошибка при получении курсов из базы данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []bson.M
	if err = cursor.All(context.Background(), &users); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	for i := range users {
		users[i]["_id"] = users[i]["_id"].(primitive.ObjectID).Hex()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
