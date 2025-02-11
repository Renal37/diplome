package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

func DownloadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Invalid courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": courseId}},
		bson.M{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "courses",
			"localField":   "courseId",
			"foreignField": "_id",
			"as":           "course",
		}},
		bson.M{"$project": bson.M{
			"user":   1,
			"course": 1,
			"status": 1,
		}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("Aggregation error: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil || len(result) == 0 {
		log.Println("No data found")
		http.Error(w, "Данные не найдены", http.StatusNotFound)
		return
	}

	// Обработка данных пользователя
	userArray, ok := result[0]["user"].(primitive.A)
	if !ok || len(userArray) == 0 {
		log.Println("User data not found or invalid format")
		http.Error(w, "Данные пользователя не найдены", http.StatusNotFound)
		return
	}
	user := userArray[0].(primitive.M)

	// Обработка данных курса
	courseArray, ok := result[0]["course"].(primitive.A)
	if !ok || len(courseArray) == 0 {
		log.Println("Course data not found or invalid format")
		http.Error(w, "Данные курса не найдены", http.StatusNotFound)
		return
	}
	course := courseArray[0].(primitive.M)

	// Создание PDF-документа
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Путь к файлам шрифтов
	fontPathRegular := "../server/font/timesnewromanpsmt.ttf"
	fontPathBold := "../server/font/timesnewromanbold.ttf"

	// Проверка существования файлов шрифтов
	_, err = os.Stat(fontPathRegular)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Regular font file not found: %v", err)
			http.Error(w, "Шрифт (Regular) не найден", http.StatusInternalServerError)
			return
		}
		log.Printf("Error checking regular font file: %v", err)
		http.Error(w, "Ошибка при проверке шрифта (Regular)", http.StatusInternalServerError)
		return
	}

	_, err = os.Stat(fontPathBold)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Bold font file not found: %v", err)
			http.Error(w, "Шрифт (Bold) не найден", http.StatusInternalServerError)
			return
		}
		log.Printf("Error checking bold font file: %v", err)
		http.Error(w, "Ошибка при проверке шрифта (Bold)", http.StatusInternalServerError)
		return
	}

	// Добавление шрифтов
	pdf.AddUTF8Font("TimesNewRoman", "", fontPathRegular)
	pdf.AddUTF8Font("TimesNewRoman", "B", fontPathBold)

	// Установка шрифта (Bold) для заголовка
	pdf.SetFont("TimesNewRoman", "B", 16)
	pdf.Cell(0, 10, "ДОГОВОР № _______")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.Cell(0, 10, "об образовании на обучение по дополнительным образовательным программам")
	pdf.Ln(10)
	pdf.Cell(0, 10, "г. Альметьевск                                  «___» ____________ 20__ г.")
	pdf.Ln(20)

	// Раздел I. Предмет Договора
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.Cell(0, 10, "I. Предмет Договора")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.MultiCell(0, 5, "1.1. Исполнитель обязуется предоставить образовательную услугу по обучению по программе, указанной в п.1.2. настоящего договора, в пределах федерального государственного образовательного стандарта и (или) профессиональных стандартов в соответствии с учебным планом, в том числе индивидуальным, и образовательной программой Исполнителя, а Обучающийся обязуется оплатить указанную образовательную услугу.", "0", "L", false)
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("1.2. Наименование дополнительной образовательной программы: %s", course["title"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("1.3. Срок обучения составляет %d часов.", course["duration"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, "1.4. Форма обучения – очно-заочная.")
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "1.5. После освоения Обучающимся образовательной программы, успешного прохождения итоговой аттестации и полностью оплатившему за обучение, в соответствии с условиями договора выдается документ установленного образца - диплом о профессиональной переподготовке.", "0", "L", false)
	pdf.Ln(20)

	// Раздел II. Права Исполнителя и Обучающегося
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.Cell(0, 10, "II. Права Исполнителя и Обучающегося")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.MultiCell(0, 5, "2.1. Исполнитель вправе:", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "2.1.1. Самостоятельно осуществлять образовательный процесс, устанавливать системы оценок, формы, порядок и периодичность проведения промежуточной аттестации Обучающегося.", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "2.1.2. Применять к Обучающемуся меры поощрения и меры дисциплинарного взыскания в соответствии с законодательством Российской Федерации, учредительными документами Исполнителя, настоящим Договором и локальными нормативными актами Исполнителя.", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "2.2. Обучающийся вправе получать информацию от Исполнителя по вопросам организации и обеспечения надлежащего предоставления услуг, предусмотренных разделом I настоящего Договора.", "0", "L", false)
	pdf.Ln(20)

	// Раздел III. Обязанности Исполнителя и Обучающегося
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.Cell(0, 10, "III. Обязанности Исполнителя и Обучающегося")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.MultiCell(0, 5, "3.1. Исполнитель обязан:", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "3.1.1. Зачислить Обучающегося, выполнившего установленные законодательством Российской Федерации, учредительными документами, локальными нормативными актами Исполнителя условия приема (зачисления).", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "3.1.2. Довести до Обучающегося информацию, содержащую сведения о предоставлении платных образовательных услуг в порядке и объеме, которые предусмотрены Законом Российской Федерации “О защите прав потребителей” и Федеральным законом “Об образовании в Российской Федерации”.", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "3.2. Обучающийся обязан:", "0", "L", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 5, "3.2.1. Выполнять задания для подготовки к занятиям, предусмотренным учебным планом, в том числе индивидуальным.", "0", "L", false)
	pdf.Ln(20)

	// Раздел IX. Адреса и реквизиты сторон
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.Cell(0, 10, "IX. Адреса и реквизиты сторон")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.Cell(0, 10, "Исполнитель:")
	pdf.Ln(5)
	pdf.Cell(0, 10, "ГАПОУ «Альметьевский политехнический техникум»")
	pdf.Ln(5)
	pdf.Cell(0, 10, "423457, РТ, г. Альметьевск, ул. Мира д.10.")
	pdf.Ln(5)
	pdf.Cell(0, 10, "ИНН/КПП 1644005722/164401001")
	pdf.Ln(5)
	pdf.Cell(0, 10, "ОГРН 1021601625352")
	pdf.Ln(10)
	pdf.Cell(0, 10, "Обучающийся:")
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("ФИО: %s %s %s", user["lastname"], user["firstname"], user["middlename"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Адрес: %s", user["homeAddress"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Паспорт: %s", user["passportData"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Телефон: %s", user["phone"]))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("E-mail: %s", user["email"]))

	// Отправка PDF-файла клиенту
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract.pdf")
	if err := pdf.Output(w); err != nil {
		log.Printf("PDF generation error: %v", err)
		http.Error(w, "Ошибка при генерации PDF", http.StatusInternalServerError)
	}
}
