import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './course_component.css';

const CourseRegistration = () => {
    const { courseId } = useParams();
    const navigate = useNavigate(); // Используем useNavigate вместо useHistory
    const [course, setCourse] = useState(null);
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        // Загружаем данные курса
        fetch(`http://localhost:5000/courses/${courseId}`)
            .then(response => response.json())
            .then(data => {
                setCourse(data);
                setIsLoading(false);
            })
            .catch(error => {
                console.error('Error fetching course:', error);
                setError('Ошибка при загрузке курса');
                setIsLoading(false);
            });

        // Загружаем данные пользователя
        fetch('http://localhost:5000/profile', {
            method: "GET",
            credentials: "include",
        })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    navigate('/auth/login');
                } else {
                    setUser(data);
                }
            })
            .catch(error => {
                console.error('Error fetching user data:', error);
                // navigate('/auth/login'); // Используем navigate вместо history.push
            });
    }, [courseId, navigate]);

    const handleRegister = () => {
        if (!user || !user.ID) {
            console.error("User ID is missing");
            return;
        }
        if (!user.lastName || !user.firstName || !user.middleName) {
            navigate('/auth/edit_profile'); // Используем navigate вместо history.push
            return;
        }

        fetch('http://localhost:5000/courses/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({
                courseId: courseId,
                userId: user.ID,
            }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Вы успешно записаны на курс! Ожидайте одобрения администратора.');
                    navigate('/'); // Используем navigate вместо history.push
                } else {
                    setError(data.message || 'Ошибка при записи на курс');
                }
            })
            .catch(error => {
                console.error('Error registering for course:', error);
                setError('Ошибка при записи на курс');
            });
    };

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div className="course-registration">
            <div className="course-registration-container">
                <h1 className="course-registration-title">Запись на курс: {course.title}</h1>
                <p className="course-registration-description">{course.description}</p>
                <div className="course-registration-details">
                    <p><strong>Продолжительность:</strong> {course.duration} часов</p>
                    <p><strong>Стоимость:</strong> {course.price} руб.</p>
                    <p><strong>Тип:</strong> {course.type}</p>
                </div>
                <button className="course-registration-button" onClick={handleRegister}>
                    Записаться на курс
                </button>
                {error && <p className="course-registration-error">{error}</p>}
            </div>
        </div>
    );
};

export default CourseRegistration;