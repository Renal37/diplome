import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import './main.center.component.css';

const Main_center = () => {
    const [courses, setCourses] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [filterType, setFilterType] = useState('');
    const location = useLocation(); // Хук для получения текущего пути

    useEffect(() => {
        fetch('http://localhost:5000/courses')
            .then(response => response.json())
            .then(data => setCourses(data))
            .catch(error => console.error('Error fetching courses:', error));
    }, []);

    // Определяем тип курсов в зависимости от текущего пути
    const getCourseTypeByPath = () => {
        if (location.pathname === '/promotion') {
            return 'Повышение квалификации';
        } else if (location.pathname === '/professional') {
            return 'Профессиональная переподготовка';
        } else {
            return ''; // На других страницах показываем все курсы
        }
    };

    // Фильтруем курсы по поисковому запросу, типу фильтра и текущему пути
    const filteredCourses = courses.filter(course => {
        const matchesSearchTerm = course.title.toLowerCase().includes(searchTerm.toLowerCase());
        const matchesFilterType = filterType ? course.type === filterType : true;
        const matchesPathType = getCourseTypeByPath() ? course.type === getCourseTypeByPath() : true;

        return matchesSearchTerm && matchesFilterType && matchesPathType;
    });

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
                    {location.pathname !== '/promotion' && location.pathname !== '/professional' && (
                        <select onChange={e => setFilterType(e.target.value)} value={filterType}>
                            <option value="">Все типы</option>
                            <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
                            <option value="Повышение квалификации">Повышение квалификации</option>
                        </select>
                    )}

                </form>
                <div className="courses_list">
                    {filteredCourses.map((course, index) => (
                        <div
                            key={course._id}
                            className={`course_item ${index % 2 === 0 ? 'even' : 'odd'}`} // Добавляем классы even и odd
                        >

                            <div className='course_text'>
                                <h2>{course.title}</h2>
                                <p>{course.description}</p>
                             
                            </div>

                            <div className='course_btn'>
                                <Link to={`/courses/register/${course._id}`}>Записаться на курс</Link>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default Main_center;