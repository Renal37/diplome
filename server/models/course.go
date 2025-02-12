package models

type Course struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // Продолжительность в часах
	Price       int    `json:"price"`    // Стоимость курса
	Type        string `json:"type"`     // Тип курса
}