import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

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
            .then(response => response.json(),
            // console.log('response:', data)
        )
            .then(data => {
                if (data.error) {
                    history.push('/auth/login');
                } else {
                    setUser(data);
                }
            })
            .catch(error => {
                console.error('Error fetching user data:', error);
                history.push('/auth/login');
            });
        }, [courseId]);

    const handleRegister = () => {
        if (!user || !user.fullName || !user.email || !user.phone) {
            history.push('/auth/edit_profile');
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
                userId: user._id,
            }),
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    alert('Вы успешно записаны на курс! Ожидайте одобрения администратора.');
                    history.push('/');
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
            <h1>Запись на курс: {course.title}</h1>
            <p>{course.description}</p>
            <p>Продолжительность: {course.duration} часов</p>
            <p>Стоимость: {course.price} руб.</p>
            <p>Тип: {course.type}</p>
            <button onClick={handleRegister}>Записаться на курс</button>
        </div>
    );
};

export default CourseRegistration;