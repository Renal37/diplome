package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Course представляет структуру данных курса
type Course struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // Продолжительность в часах
	Price       int    `json:"price"`    // Стоимость курса
	Type        string `json:"type"`     // Тип курса
}

// User представляет структуру данных пользователя
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FullName     string             `json:"fullName"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	BirthDate    string             `json:"birthDate"`
	Residence    string             `json:"residence,omitempty"`
	Education    string             `json:"education,omitempty"`
	BirthPlace   string             `json:"birthPlace,omitempty"`
	HomeAddress  string             `json:"homeAddress,omitempty"`
	PassportData string             `json:"passportData,omitempty"`
	SNILS        string             `json:"snils,omitempty"`
}

// JWT ключ для подписи токенов
var jwtKey = []byte("your_secret_key")

// Claims представляет структуру данных для токена
type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Change to your client URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func main() {
	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Создание маршрутов
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/add-course", addCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/courses", getCourses).Methods("GET")
	r.HandleFunc("/update-course/{id}", updateCourse).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete-course/{id}", deleteCourse).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/register", registerUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-profile/{id}", updateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/check-token", checkToken).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", loginUser).Methods("POST", "OPTIONS")

	// Запуск сервера
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

// addCourse обрабатывает запросы на добавление нового курса
func addCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Подключение к коллекции courses
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	// Вставка данных курса в коллекцию
	_, err = collection.InsertOne(context.Background(), course)
	if err != nil {
		http.Error(w, "Ошибка при добавлении курса в базу данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс '%s' успешно добавлен!", course.Title)
}

// getCourses обрабатывает запросы на получение списка курсов
func getCourses(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Ошибка при получении курсов из базы данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	// Преобразуем ObjectId в строку для передачи в JSON
	for i := range courses {
		courses[i]["_id"] = courses[i]["_id"].(primitive.ObjectID).Hex()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// updateCourse обрабатывает запросы на обновление курса по его идентификатору
func updateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var course Course
	err = json.NewDecoder(r.Body).Decode(&course)
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
	collection := client.Database("diplome").Collection("courses")

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       course.Title,
			"description": course.Description,
			"duration":    course.Duration,
			"price":       course.Price,
			"type":        course.Type,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении курса в базе данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс '%s' успешно обновлен!", course.Title)
}

// deleteCourse обрабатывает запросы на удаление курса по его идентификатору
func deleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении курса из базы данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс успешно удален!")
}

// registerUser обрабатывает запросы на регистрацию нового пользователя
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при хешировании пароля", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Подключение к коллекции users
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("users")

	// Вставка данных пользователя в коллекцию
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
		return
	}

	// Создание JWT токена
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: result.InsertedID.(primitive.ObjectID).Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Отправка токена клиенту
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	fmt.Fprintf(w, "Пользователь '%s' успешно зарегистрирован!", user.FullName)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
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

	var user User
	err = collection.FindOne(context.Background(), bson.M{"username": credentials.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   false, // Set to true if using HTTPS
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// updateProfile обрабатывает запросы на обновление профиля пользователя
func updateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
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

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"residence":    user.Residence,
			"education":    user.Education,
			"birthPlace":   user.BirthPlace,
			"homeAddress":  user.HomeAddress,
			"passportData": user.PassportData,
			"snils":        user.SNILS,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении профиля в базе данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Профиль пользователя успешно обновлен!")
}

// checkToken обрабатывает запросы на проверку токена
func checkToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	tokenStr = tokenStr[len("Bearer "):]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		fmt.Println("Неверный токен:", err)
		return
	}

	fmt.Println("Токен действителен для пользователя:", claims.UserID)
	w.WriteHeader(http.StatusOK)
}

func checkCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Cookie не найдено
			fmt.Println("Cookie не найдено")
			http.Error(w, "Cookie не найдено", http.StatusUnauthorized)
			return
		}
		// Другая ошибка
		fmt.Println("Ошибка при чтении cookie:", err)
		http.Error(w, "Ошибка при чтении cookie", http.StatusBadRequest)
		return
	}

	// Cookie найдено
	fmt.Println("Cookie найдено:", cookie.Value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": cookie.Value})
}
