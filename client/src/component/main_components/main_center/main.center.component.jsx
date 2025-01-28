import React, { useState, useEffect } from 'react';
import './main.center.component.css';

const Main_center = () => {
    const [courses, setCourses] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [filterType, setFilterType] = useState('');

    useEffect(() => {
        fetch('http://localhost:5000/courses')
            .then(response => response.json())
            .then(data => setCourses(data))
            .catch(error => console.error('Error fetching courses:', error));
    }, []);

    const filteredCourses = courses.filter(course =>
        course.title.toLowerCase().includes(searchTerm.toLowerCase()) &&
        (filterType ? course.type === filterType : true)
    );

    return (
        <div className="main_center">
            <div className="main_center_header">
                <div className="line"></div>
                <h1>Поиск по всем программам</h1>
                <div className="line"></div>
            </div>
            <div className="main_center_center">
                <form onSubmit={e => e.preventDefault()}>
                    <input
                        type="text"
                        placeholder="Название программы обучения"
                        value={searchTerm}
                        onChange={e => setSearchTerm(e.target.value)}
                    />
                    <select onChange={e => setFilterType(e.target.value)} value={filterType}>
                        <option value="">Все типы</option>
                        <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
                        <option value="Повышение квалификации">Повышение квалификации</option>
                    </select>
                </form>
                <div className="courses_list">
                    {filteredCourses.map(course => (
                        <div key={course._id} className="course_item">
                            <h2>{course.title}</h2>
                            <p>{course.description}</p>
                            <p>Продолжительность: {course.duration} часов</p>
                            <p>Стоимость: {course.price} руб.</p>
                            <p>Тип: {course.type}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default Main_center;