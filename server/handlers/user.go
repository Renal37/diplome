package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.UsersCollection)

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
    // Получаем токен из куки
    cookie, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
            return
        }
        http.Error(w, "Ошибка при получении токена", http.StatusBadRequest)
        return
    }

    // Парсим токен
    tokenStr := cookie.Value
    claims := &utils.Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return utils.JwtKey, nil
    })

    if err != nil || !token.Valid {
        http.Error(w, "Неверный токен", http.StatusUnauthorized)
        return
    }

    // Преобразуем ID пользователя в ObjectID
    userID, err := primitive.ObjectIDFromHex(claims.UserID)
    if err != nil {
        http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
        return
    }

    // Получаем коллекцию пользователей
    collection := db.GetCollection(db.UsersCollection)

    // Создаем pipeline для агрегации
    pipeline := []bson.M{
        {
            "$match": bson.M{"_id": userID},
        },
        {
            "$lookup": bson.M{
                "from":         db.EducationsCollection,
                "localField":   "educationId",
                "foreignField": "_id",
                "as":           "education",
            },
        },
        {
            "$addFields": bson.M{
                "education": bson.M{"$arrayElemAt": bson.A{"$education", 0}},
            },
        },
        {
            "$project": bson.M{
                "password": 0, // Исключаем пароль из результатов
            },
        },
    }

    // Выполняем агрегацию
    cursor, err := collection.Aggregate(context.Background(), pipeline)
    if err != nil {
        http.Error(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    // Декодируем результат
    var userData bson.M
    if cursor.Next(context.Background()) {
        if err := cursor.Decode(&userData); err != nil {
            http.Error(w, "Ошибка при обработке данных пользователя", http.StatusInternalServerError)
            return
        }
    } else {
        http.Error(w, "Пользователь не найден", http.StatusNotFound)
        return
    }

    // Преобразуем ObjectID в строку для удобства фронтенда
    if userData["_id"] != nil {
        userData["_id"] = userData["_id"].(primitive.ObjectID).Hex()
    }
    if userData["educationId"] != nil {
        userData["educationId"] = userData["educationId"].(primitive.ObjectID).Hex()
    }
    if userData["education"] != nil {
        education := userData["education"].(bson.M)
        if education["_id"] != nil {
            education["_id"] = education["_id"].(primitive.ObjectID).Hex()
        }
    }

    // Отправляем данные пользователя
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(userData); err != nil {
        http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
    }
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
    // Получаем и парсим данные запроса
    var updateData struct {
        LastName          string `json:"lastName,omitempty"`
        FirstName         string `json:"firstName,omitempty"`
        MiddleName        string `json:"middleName,omitempty"`
        EducationID       string `json:"educationId,omitempty"` // Изменено на строку
        Phone             string `json:"phone,omitempty"`
        BirthDate         string `json:"birthDate,omitempty"`
        BirthPlace        string `json:"birthPlace,omitempty"`
        HomeAddress       string `json:"homeAddress,omitempty"`
        WorkPlace         string `json:"workPlace,omitempty"`
        JobTitle          string `json:"jobTitle,omitempty"`
        OldPassword       string `json:"oldPassword,omitempty"`
        NewPassword       string `json:"newPassword,omitempty"`
        AgreeToProcessing bool   `json:"agreeToProcessing,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
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

    userID, err := primitive.ObjectIDFromHex(claims.UserID)
    if err != nil {
        http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection(db.UsersCollection)

    // Создаём объект для обновления
    update := bson.M{}

    // Обновляем основные поля
    if updateData.LastName != "" {
        update["lastname"] = updateData.LastName
    }
    if updateData.FirstName != "" {
        update["firstname"] = updateData.FirstName
    }
    if updateData.MiddleName != "" {
        update["middlename"] = updateData.MiddleName
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
    update["agreetoprocessing"] = updateData.AgreeToProcessing

    // Обработка образования
    if updateData.EducationID != "" {
        educationID, err := primitive.ObjectIDFromHex(updateData.EducationID)
        if err != nil {
            http.Error(w, "Неверный формат ID образования", http.StatusBadRequest)
            return
        }
        update["educationid"] = educationID // Важно: используем правильное имя поля в БД
    } else {
        update["educationid"] = nil // Очищаем, если образование не указано
    }

    // Обработка пароля
    if updateData.OldPassword != "" && updateData.NewPassword != "" {
        var user models.User
        if err := collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user); err != nil {
            http.Error(w, "Пользователь не найден", http.StatusNotFound)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateData.OldPassword)); err != nil {
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

    // Выполняем обновление
    result, err := collection.UpdateOne(
        context.Background(),
        bson.M{"_id": userID},
        bson.M{"$set": update},
    )

    if err != nil {
        http.Error(w, "Ошибка при обновлении профиля", http.StatusInternalServerError)
        return
    }

    if result.MatchedCount == 0 {
        http.Error(w, "Пользователь не найден", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Профиль успешно обновлён",
    })
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.UsersCollection)

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
