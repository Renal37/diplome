import React, { useState, useEffect } from 'react';
import './admin_update_component.css';

const AdminUpdateComponent = () => {
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePrice, setCoursePrice] = useState('');
  const [courseType, setCourseType] = useState('');

  const fetchCourses = async () => {
    const response = await fetch('http://localhost:5000/courses');
    const data = await response.json();
    setCourses(data);
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  const handleCourseSelect = (course) => {
    setSelectedCourse(course);
    setCourseTitle(course.title);
    setCourseDescription(course.description);
    setCourseDuration(course.duration);
    setCoursePrice(course.price);
    setCourseType(course.type);
  };

  const handleUpdate = async (e) => {
    e.preventDefault();
    const response = await fetch(`http://localhost:5000/update-course/${selectedCourse._id}`, {
      method: 'PUT',
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
      alert('Course updated successfully');
      fetchCourses();
      setSelectedCourse(null);
      setCourseTitle('');
      setCourseDescription('');
      setCourseDuration('');
      setCoursePrice('');
      setCourseType('');
    } else {
      alert('Failed to update course');
    }
  };

  return (
    <div className="admin-update-component">
      <h1>Обновление курсов</h1>
      <table className="course-table">
        <thead>
          <tr>
            <th>Заголовок</th>
            <th>Тип</th>
            <th>Стоимость</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {courses.map((course) => (
            <tr key={course._id} onClick={() => handleCourseSelect(course)}>
              <td>{course.title}</td>
              <td>{course.type}</td>
              <td>{course.price} руб.</td>
              <td><button onClick={() => handleCourseSelect(course)}>Редактировать</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      {selectedCourse && (
        <form onSubmit={handleUpdate}>
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
          <button type="submit">Обновить курс</button>
        </form>
      )}
    </div>
  );
};

export default AdminUpdateComponent;