import { useState, useEffect } from "react";
import "./check_course_component.css";

const CheckCourse = () => {
    const [courses, setCourses] = useState([]); // Инициализируем как пустой массив
    const [loading, setLoading] = useState(true); // Состояние загрузки
    const [error, setError] = useState(null); // Состояние ошибки
    const [selectedStatus, setSelectedStatus] = useState("all"); // Текущий статус

    // Функция для получения курсов
    const fetchCourses = async (status) => {
        try {
            let endpoint = "/user/courses"; // По умолчанию все курсы пользователя
            if (status !== "all") {
                endpoint = `/user/courses/status?status=${status}`; // Курсы по статусу
            }
            const response = await fetch(`http://localhost:5000${endpoint}`, {
                method: "GET",
                credentials: "include", // Включаем аутентификацию через cookies
            });
            if (!response.ok) {
                throw new Error("Ошибка при загрузке курсов");
            }
            const data = await response.json();
            // Проверяем, что данные являются массивом
            if (!Array.isArray(data)) {
                console.error("Сервер вернул некорректные данные:", data);
                setCourses([]);
                return;
            }
            setCourses(data);
        } catch (err) {
            setError(err.message);
            setCourses([]); // При ошибке очищаем список курсов
        } finally {
            setLoading(false);
        }
    };

    // Загрузка данных при изменении selectedStatus
    useEffect(() => {
        fetchCourses(selectedStatus);

    }, [selectedStatus]);

    // Обработка отображения
    if (loading) return <p className="loading">Загрузка...</p>;
    if (error) return <p className="error">{error}</p>;
    console.log(courses);

    return (
        <div className="check-course-container">
            {/* Навигация по статусам */}
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

            {/* Отображение списка курсов */}
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
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
};

export default CheckCourse;