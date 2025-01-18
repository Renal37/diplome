import React, { useState } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');

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
      }),
    });
    if (response.ok) {
      alert('Course added successfully');
    } else {
      alert('Failed to add course');
    }
  };

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
        <button type="submit">Добавить курс</button>
      </form>
    </div>
  );
};

export default AdminCourseManagement;