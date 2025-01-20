import React, { useState, useEffect } from 'react';

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
    if (!selectedCourse) {
      alert('Пожалуйста, выберите курс для обновления');
      return;
    }
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
    <div>
      <h1>Обновление курсов</h1>
      <ul>
        {courses.map((course) => (
          <li key={course._id} onClick={() => handleCourseSelect(course)}>
            {course.title} - {course.type} - {course.price} руб.
          </li>
        ))}
      </ul>
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