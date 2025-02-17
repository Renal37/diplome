import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './course_component.css';

const CourseRegistration = () => {
    const { courseId } = useParams();
    const navigate = useNavigate();
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
                    console.log(data);
                }
            })
            .catch(error => {
                console.error('Error fetching user data:', error);
                // navigate('/auth/login'); 
            });
    }, [courseId, navigate]);

    const handleRegister = () => {
        if (!user || !user.ID) {
            console.error("User ID is missing");
            return;
        }

        // Проверка уровня образования для курсов типа "Профессиональная переподготовка"
        if (course.type === "Профессиональная переподготовка") {
            const allowedEducations = ["Среднее профессиональное", "Высшее"];
            if (!allowedEducations.includes(user.education)) {
                alert("Для записи на этот курс требуется среднее профессиональное или высшее образование.");
                return;
            }
        }

        const requiredFields = [
            { field: 'lastName', message: 'Фамилия не заполнена' },
            { field: 'firstName', message: 'Имя не заполнено' },
            { field: 'middleName', message: 'Отчество не заполнено' },
            { field: 'birthDate', message: 'Дата рождения не указана' },
            { field: 'birthPlace', message: 'Место рождения не указано' },
            { field: 'education', message: 'Образование не указано' },
            { field: 'email', message: 'Email не указан' },
            { field: 'homeAddress', message: 'Домашний адрес не указан' },
            { field: 'jobTitle', message: 'Должность не указана' },
            { field: 'passportData', message: 'Паспортные данные не указаны' },
            { field: 'phone', message: 'Телефон не указан' },
            { field: 'snils', message: 'СНИЛС не указан' },
            { field: 'workPlace', message: 'Место работы не указано' },
            { field: 'contractuploaded', message: 'Соглашение не найдено' },
        ];

        for (const { field, message } of requiredFields) {
            if (!user[field]) {
                alert(message);
                navigate("/auth/edit_profile");
                return;
            }
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
                    navigate('/');
                } else {
                    setError(data.message || 'Ошибка при записи на курс');
                }
            })
            .catch(error => {
                console.error('Error registering for course:', error);
                setError('Ошибка при записи на курс');
            });
    };

    const isProfileComplete = () => {
        return user && user.ID;
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
                <div className='course-registration-text'>
                    <h1 className="course-registration-title">Запись на курс: {course.title}</h1>
                    <p className="course-registration-description">{course.description}</p>
                </div>
                <div className="course-registration-details">
                    <p><strong>Продолжительность:</strong> {course.duration} часов</p>
                    <p><strong>Стоимость:</strong> {course.price} руб.</p>
                    <p><strong>Тип:</strong> {course.type}</p>
                </div>
                <div className='course_btn'>
                    {isProfileComplete() && (
                        <button className="course-registration-button" onClick={handleRegister}>
                            Записаться на курс
                        </button>
                    )}
                </div>

                {error && <p className="course-registration-error">{error}</p>}
            </div>
        </div>
    );
};

export default CourseRegistration;