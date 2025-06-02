import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './course_component.css';

const CourseRegistration = () => {
    const { courseId } = useParams();
    const navigate = useNavigate();
    const [course, setCourse] = useState(null);
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    const fetchCourseData = async () => {
        try {
            const response = await fetch(`http://localhost:5000/courses/${courseId}`);
            if (!response.ok) {
                throw new Error('Ошибка при загрузке курса');
            }
            const data = await response.json();
            console.log('Course data:', data);
            setCourse(data);
            setIsLoading(false);
        } catch (error) {
            console.error('Error fetching course:', error);
            alert(error.message);
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchCourseData();

        fetch('http://localhost:5000/profile', {
            method: "GET",
            credentials: "include",
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка при загрузке профиля');
                }
                return response.json();
            })
            .then(data => {
                if (data.error) {
                    navigate('/auth/login');
                } else {
                    setUser(data);
                    console.log('User data:', data);
                }
            })
            .catch(error => {
                console.error('Error fetching user data:', error);
                alert('Необходимо авторизоваться');
                navigate('/auth/login');
            });
    }, [courseId, navigate]);

    const handleRegister = async () => {
        if (!user || !user._id) {
            alert('Пользователь не авторизован');
            navigate('/auth/login');
            return;
        }

        await fetchCourseData();

        const currentDate = new Date();
        const registrationStart = new Date(course.registrationStart);
        const registrationEnd = new Date(course.registrationEnd);

        if (currentDate < registrationStart) {
            alert('Регистрация на курс ещё не началась');
            return;
        }
        if (currentDate > registrationEnd) {
            alert('Регистрация на курс уже закончилась');
            return;
        }

        if (course.activeStudentsCount >= course.maxStudents) {
            alert(`Курс достиг максимального количества студентов (${course.activeStudentsCount}/${course.maxStudents}). Свободных мест нет.`);
            return;
        }

        if (course.type === "Профессиональная переподготовка") {
            const allowedEducations = [
                "Среднее профессиональное",
                "Высшее",
                "Высшее образование",
            ];
            if (!user.education || !allowedEducations.includes(user.education.name)) {
                alert(
                    'Для записи на этот курс требуется среднее профессиональное или высшее образование.'
                );
                return;
            }
        }

        const requiredFields = [
            { field: 'lastname', message: 'Фамилия не заполнена' },
            { field: 'firstname', message: 'Имя не заполнено' },
            { field: 'middlename', message: 'Отчество не заполнено' },
            { field: 'birthdate', message: 'Дата рождения не указана' },
            { field: 'birthplace', message: 'Место рождения не указано' },
            { field: 'education', message: 'Образование не указано' },
            { field: 'email', message: '' },
            { field: 'homeaddress', message: 'Домашний адрес не указан' },
            { field: 'jobtitle', message: 'Должность не указана' },
            { field: 'passportdata', message: 'Паспортные данные не указаны' },
            { field: 'phone', message: 'Телефон не указан' },
            { field: 'snils', message: 'СНИЛС не указан' },
            { field: 'workplace', message: 'Место работы не указано' },
            { field: 'passportissuedby', message: 'Кем выдан паспорт не указано' },
            { field: 'passportissuedate', message: 'Дата выдачи паспорта не указана' },
            { field: 'agreetoprocessing', message: 'Согласие на обработку данных не получено' },
            { field: 'contractUploaded', message: 'Соглашение не загружено' },
        ];

        for (const { field, message } of requiredFields) {
            if (field === 'education') {
                if (!user[field] || !user[field].name) {
                    alert(message);
                    navigate('/auth/edit_profile');
                    return;
                }
            } else if (field === 'agreetoprocessing' || field === 'contract') {
                if (!user[field]) {
                    alert(message);
                    return;
                }
            } else if (!user[field]) {
                alert(message);
                navigate('/auth/edit_profile');
                return;
            }
        }

        try {
            const response = await fetch('http://localhost:5000/courses/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    courseId: courseId,
                    userId: user._id,
                }),
            });

            const data = await response.json();
            if (!response.ok) {
                throw new Error(data.error || 'Ошибка при записи на курс');
            }

            if (data.success) {
                alert('Вы успешно записаны на курс! Ожидайте одобрения администратора.');
                navigate('/');
            } else {
                alert(data.error || 'Ошибка при записи на курс');
            }
        } catch (error) {
            console.error('Error registering for course:', error);
            alert(error.message);
        }
    };

    const isProfileComplete = () => {
        return user && user._id;
    };

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="course-registration">
            <div className="course-registration-container">
                <div className="course-registration-text">
                    <h1 className="course-registration-title">Запись на курс: {course.title}</h1>
                    <p className="course-registration-description">{course.description}</p>
                </div>
                <div className="course-registration-details">
                    <p><strong>Продолжительность:</strong> {course.duration} часов</p>
                    <p><strong>Стоимость:</strong> {course.price} руб.</p>
                    <p><strong>Тип:</strong> {course.type}</p>
                    <p><strong>Начало регистрации:</strong> {new Date(course.registrationStart).toLocaleDateString()}</p>
                    <p><strong>Окончание регистрации:</strong> {new Date(course.registrationEnd).toLocaleDateString()}</p>
                    <p><strong>Осталось мест:</strong> {course.maxStudents - (course.activeStudentsCount || 0)}</p>
                </div>
                <div className="course_btn">
                    {isProfileComplete() && (
                        <button className="course-registration-button" onClick={handleRegister}>
                            Записаться на курс
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default CourseRegistration;