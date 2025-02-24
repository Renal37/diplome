import React, { useState, useEffect } from 'react';
import './admin_add_component.css';

const AdminAddComponent = () => {
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePrice, setCoursePrice] = useState('');
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
        price: parseInt(coursePrice, 10),
        type: courseType,
      }),
    });
    if (response.ok) {
      alert('Course added successfully');
      // Clear form fields
      setCourseTitle('');
      setCourseDescription('');
      setCourseDuration('');
      setCoursePrice('');
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
      <h1>Управление курсами</h1>
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
          placeholder="Продолжительность (часы)"
          value={courseDuration}
          onChange={(e) => setCourseDuration(e.target.value)}
        />
        <input
          type="number"
          placeholder="Стоимость"
          value={coursePrice}
          onChange={(e) => setCoursePrice(e.target.value)}
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
    </div>
  );
};

export default AdminAddComponent;