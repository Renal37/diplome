import { useState, useEffect } from "react";
import "./check_course_component.css";
import { PDFDocument } from "pdf-lib"; // Импортируем PDF-Lib

const CheckCourse = () => {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedStatus, setSelectedStatus] = useState("all");
    const [userData, setUserData] = useState({}); // Добавляем состояние для данных пользователя

    const fetchCourses = async (status) => {
        try {
            let endpoint = "/user/courses";
            if (status !== "all") {
                endpoint = `/user/courses/status?status=${status}`;
            }
            const response = await fetch(`http://localhost:5000${endpoint}`, {
                method: "GET",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при загрузке курсов");
            }
            const data = await response.json();
            if (!Array.isArray(data)) {
                console.error("Сервер вернул некорректные данные:", data);
                setCourses([]);
                return;
            }
            setCourses(data);
        } catch (err) {
            setError(err.message);
            setCourses([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCourses(selectedStatus);
    }, [selectedStatus]);

    // Функция для заполнения и скачивания PDF-договора
    const handleDownloadContract = async (courseId) => {
        try {
            // Загружаем PDF-шаблон с сервера
            const response = await fetch(`http://localhost:5000/user/download-contract/${courseId}`, {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Ошибка при скачивании договора");
            }

            const pdfBytes = await response.arrayBuffer(); // Получаем байты PDF
            const pdfDoc = await PDFDocument.load(pdfBytes); // Загружаем PDF в PDF-Lib

            // Пример данных для заполнения полей (замените на реальные данные пользователя)
            const dataToFill = {
                // Text1:  "Иванов Иван Иванович",
                // Address: userData.address || "г. Москва, ул. Ленина, д. 1",
                // Phone: userData.phone || "+7 (999) 123-45-67",
                // Email: userData.email || "example@example.com",
                // CourseName: courses.find((course) => course._id === courseId)?.courseTitle || "Название курса",
            };

            // Заполняем поля формы
            for (const [fieldName, fieldValue] of Object.entries(dataToFill)) {
                const field = form.getField(fieldName);
                if (field instanceof TextField) {
                    field.setText(fieldValue);
                } else {
                    console.warn(`Поле '${fieldName}' не является текстовым.`);
                }
            }

            // Сохраняем изменения
            const updatedPdfBytes = await pdfDoc.save();

            // Создаем ссылку для скачивания
            const blob = new Blob([updatedPdfBytes], { type: "application/pdf" });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            a.download = "contract.pdf";
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            alert("Вы скачали договор, теперь заполните его и отправте нам!");
        } catch (err) {
            alert(err.message);
        }
    };

    if (loading) return <p className="loading">Загрузка...</p>;
    if (error) return <p className="error">{error}</p>;

    return (
        <div className="check-course-container">
            {/* Навигация по статусам курсов */}
            <div className="profile_course_nav">
                <button
                    onClick={() => setSelectedStatus("all")}
                    className={`prifle_nav_button ${selectedStatus === "all" ? "active" : ""}`}
                >
                    Все курсы
                </button>
                <button
                    onClick={() => setSelectedStatus("Одобренный")}
                    className={`prifle_nav_button ${selectedStatus === "Одобренный" ? "active" : ""}`}
                >
                    Одобренные
                </button>
                <button
                    onClick={() => setSelectedStatus("Отклоненный")}
                    className={`prifle_nav_button ${selectedStatus === "Отклоненный" ? "active" : ""}`}
                >
                    Отклоненные
                </button>
                <button
                    onClick={() => setSelectedStatus("Ожидание")}
                    className={`prifle_nav_button ${selectedStatus === "Ожидание" ? "active" : ""}`}
                >
                    Ожидают одобрения
                </button>
            </div>

            {/* Список курсов */}
            <div className="course-list-container">
                <h2 className="course-list-title">Список курсов:</h2>
                {courses.length === 0 ? (
                    <p className="no-courses-message">Нет доступных курсов для выбранного статуса.</p>
                ) : (
                    <ul className="course-list">
                        {courses.map((course, index) => (
                            <li key={index} className="course-item">
                                <span className="course-title">Название: {course.courseTitle}</span>
                                <span className="course-status">Статус: {course.status}</span>
                                {course.status === "Отклоненный" && course.rejectReason && (
                                    <span className="reject-reason">
                                        Причина отказа: {course.rejectReason}
                                    </span>
                                )}
                                {course.status === "Одобренный" && (
                                    <button
                                        className="download-contract-button"
                                        onClick={() => handleDownloadContract(course._id)}
                                    >
                                        Скачать договор
                                    </button>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            {/* Форма для ввода данных пользователя */}
            
        </div>
    );
};

export default CheckCourse;