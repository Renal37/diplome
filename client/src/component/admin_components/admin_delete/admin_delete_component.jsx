import React, { useState, useEffect } from 'react';
import './admin_delete_component.css';

const AdminDeleteComponent = () => {
  const [courses, setCourses] = useState([]);

  const fetchCourses = async () => {
    const response = await fetch('http://localhost:5000/courses');
    const data = await response.json();
    setCourses(data);
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  const handleDelete = async (courseId) => {
    const response = await fetch(`http://localhost:5000/delete-course/${courseId}`, {
      method: 'DELETE',
    });
    if (response.ok) {
      alert('Курс успешно удален!');
      fetchCourses();
    } else {
      alert('Ошибка при удалении курса');
    }
  };

  return (
    <div className="admin-delete-component">
      <h1>Удаление курсов</h1>
      <table className="course-table">
        <thead>
          <tr>
            <th>Заголовок</th>
            <th>Описание</th>
            <th>Продолжительность (часы)</th>
            <th>Стоимость (руб.)</th>
            <th>Тип</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {courses.map((course) => (
            <tr key={course._id}>
              <td>{course.title}</td>
              <td className="course-description">{course.description}</td>
              <td>{course.duration}</td>
              <td>{course.price}</td>
              <td>{course.type}</td>
              <td>
                <button onClick={() => handleDelete(course._id)}>Удалить</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminDeleteComponent;