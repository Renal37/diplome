import { useState, useEffect } from "react";
import './check_course_component.css';

const CheckCourse = () => {
    const [courses, setCourses] = useState([]); // Список курсов
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
                credentials: "include", // Включаем аутентификацию через токен
            });

            if (!response.ok) {
                throw new Error("Ошибка при загрузке курсов");
            }

            const data = await response.json();
            setCourses(data);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Загрузка данных при изменении selectedStatus
    useEffect(() => {
        fetchCourses(selectedStatus);
    }, [selectedStatus]);

    // Обработка отображения
    if (loading) return <p>Загрузка...</p>;
    if (error) return <p>{error}</p>;

    return (
        <div>
            {/* Навигация по статусам */}
            <div className="profile_course_nav">
                <button
                    onClick={() => setSelectedStatus("all")}
                    className="prifle_nav_button"
                >
                    Все курсы
                </button>
                <button
                    onClick={() => setSelectedStatus("Одобренный")}
                    className="prifle_nav_button"
                >
                    Одобренные
                </button>
                <button
                    onClick={() => setSelectedStatus("Отклоненный")}
                    className="prifle_nav_button"
                >
                    Отклоненные
                </button>
                <button
                    onClick={() => setSelectedStatus("Ожидание")}
                    className="prifle_nav_button"
                >
                    Ожидают одобрения
                </button>
            </div>

            {/* Отображение списка курсов */}
            <div>
                <h2>Список курсов:</h2>
                {courses.length === 0 ? (
                    <p>Нет доступных курсов для выбранного статуса.</p>
                ) : (
                    <ul>
                        {courses.map((course, index) => (
                            <li key={index}>
                                Название: {course.courseTitle}, Статус: {course.status}
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
};

export default CheckCourse;