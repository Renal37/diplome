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
    if (!course || !course._id) {
      console.error("Ошибка: отсутствует ID курса");
      return;
    }
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

  const handleClose = () => {
    setSelectedCourse(null);
    setCourseTitle('');
    setCourseDescription('');
    setCourseDuration('');
    setCoursePrice('');
    setCourseType('');
  };

  return (
    <div className="admin-update-component">
      <h1>Обновление курсов</h1>
      <table className="course-table">
        <thead>
          <tr>
            <th>Заголовок</th>
            <th>Описание</th>
            <th>Продолжительность (часы)</th>
            <th>Стоимость (руб.)</th>
            <th>Тип</th>
          </tr>
        </thead>
        <tbody>
          {courses.map((course) => (
            <React.Fragment key={course._id}>
              <tr onClick={() => handleCourseSelect(course)}>
                <td>{course.title}</td>
                <td className="course-description">{course.description}</td>
                <td>{course.duration}</td>
                <td>{course.price}</td>
                <td>{course.type}</td>
              </tr>
              {selectedCourse && selectedCourse._id === course._id && (
                <tr>
                  <td colSpan="5">
                    <form onSubmit={handleUpdate} className="update-form">
                      <input
                        type="text"
                        placeholder="Заголовок"
                        value={courseTitle}
                        onChange={(e) => setCourseTitle(e.target.value)}
                      />
                      <textarea
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
                      <div className="form-buttons">
                        <button type="submit">Обновить курс</button>
                        <button type="button" onClick={handleClose}>Закрыть</button>
                      </div>
                    </form>
                  </td>
                </tr>
              )}
            </React.Fragment>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminUpdateComponent;