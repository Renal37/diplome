import React, { useState, useEffect } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [courseType, setCourseType] = useState('Повышение квалификации');
  const [courses, setCourses] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch('http://localhost:5000/add-course', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        title: courseTitle,
        description: courseDescription,
        duration: parseInt(courseDuration, 10),
        type: courseType,
      }),
    });
    if (response.ok) {
      alert('Course added successfully');
      // Clear form fields
      setCourseTitle('');
      setCourseDescription('');
      setCourseDuration('');
      setCourseType('Повышение квалификации');
      fetchCourses();
    } else {
      alert('Failed to add course');
    }
  };

  const fetchCourses = async () => {
    const response = await fetch('http://localhost:5000/courses');
    const data = await response.json();
    setCourses(data);
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  return (
    <div className="admin-course-management">
      <h1>Hello admin!</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Заголовок"
          value={courseTitle}
          onChange={(e) => setCourseTitle(e.target.value)}
        />
        <input
          type="text"
          placeholder="Описание"
          value={courseDescription}
          onChange={(e) => setCourseDescription(e.target.value)}
        />
        <input
          type="number"
          placeholder="Стоимость"
          value={courseDuration}
          onChange={(e) => setCourseDuration(e.target.value)}
        />
        <select
          value={courseType}
          onChange={(e) => setCourseType(e.target.value)}
        >
          <option value="Повышение квалификации">Повышение квалификации</option>
          <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
        </select>
        <button type="submit">Добавить курс</button>
      </form>
      <h2>Добавленные курсы</h2>
      <ul>
        {courses.map((course, index) => (
          <li key={index}>
            {course.title} - {course.type}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default AdminCourseManagement;