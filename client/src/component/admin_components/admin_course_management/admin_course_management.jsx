import React, { useState, useEffect } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  const [courses, setCourses] = useState([]); // Инициализируем как пустой массив
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePrice, setCoursePrice] = useState('');
  const [courseType, setCourseType] = useState('Повышение квалификации');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [filterType, setFilterType] = useState('Все'); // Состояние для фильтрации по типу курса

  const fetchCourses = async () => {
    try {
      const response = await fetch('http://localhost:5000/courses');
      const data = await response.json();
      setCourses(data || []); // Убедимся, что data не null или undefined
    } catch (error) {
      console.error('Ошибка при загрузке курсов:', error);
      setCourses([]); // В случае ошибки устанавливаем пустой массив
    }
  };

  useEffect(() => {
    fetchCourses();
  }, []);

  // Обработчик нажатия клавиши Escape
  useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.key === 'Escape') {
        if (isModalOpen) {
          setIsModalOpen(false); // Закрыть модальное окно добавления курса
        }
        if (selectedCourse) {
          handleCloseUpdateForm(); // Закрыть форму обновления курса
        }
      }
    };

    // Добавляем обработчик события keydown
    window.addEventListener('keydown', handleKeyDown);

    // Убираем обработчик при размонтировании компонента
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [isModalOpen, selectedCourse]); // Зависимости от состояний модального окна и выбранного курса

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

  const handleAddCourse = async (e) => {
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
      alert('Курс был успешно добавлен');
      setIsModalOpen(false);
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

  const handleUpdateCourse = async (e) => {
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
      alert('Курс был успешно изменен');
      fetchCourses();
      setSelectedCourse(null);
      setCourseTitle('');
      setCourseDescription('');
      setCourseDuration('');
      setCoursePrice('');
      setCourseType('');
    } else {
      alert('Не удалось обновить курс');
    }
  };

  const handleDeleteCourse = async (courseId) => {
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

  const handleCloseUpdateForm = () => {
    setSelectedCourse(null);
    setCourseTitle('');
    setCourseDescription('');
    setCourseDuration('');
    setCoursePrice('');
    setCourseType('');
  };

  // Фильтрация курсов по типу
  const filteredCourses = filterType === 'Все'
    ? courses
    : courses.filter(course => course.type === filterType);

  return (
    <div className="admin-course-management">
      <div className="admin-course-header">
        {/* Фильтр по типу курса */}
        <div className="filter-section">
          <label>Фильтр по типу курса:</label>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
          >
            <option value="Все">Все</option>
            <option value="Повышение квалификации">Повышение квалификации</option>
            <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
          </select>
        </div>

        <button className='approve-btn' onClick={() => setIsModalOpen(true)}>Добавить курс</button>
      </div>
      {isModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Добавить курс</h2>
            <form onSubmit={handleAddCourse}>
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
              <div className="form-buttons">
                <button className='reject-btn' type="button" onClick={() => setIsModalOpen(false)}>Закрыть</button>
                <button className='approve-btn' type="submit">Добавить</button>
              </div>
            </form>
          </div>
        </div>
      )}

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
          {filteredCourses.map((course) => (
            <React.Fragment key={course._id}>
              <tr onClick={() => handleCourseSelect(course)}>
                <td>{course.title}</td>
                <td className="course-description">{course.description}</td>
                <td>{course.duration}</td>
                <td>{course.price}</td>
                <td>{course.type}</td>
                <td>
                  <button className='reject-btn' onClick={() => handleDeleteCourse(course._id)}>Удалить</button>
                </td>
              </tr>
              {selectedCourse && selectedCourse._id === course._id && (
                <tr>
                  <td colSpan="6">
                    <form onSubmit={handleUpdateCourse} className="update-form">
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
                        <button className='approve-btn' type="submit">Обновить курс</button>
                        <button className='reject-btn' type="button" onClick={handleCloseUpdateForm}>Закрыть</button>
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

export default AdminCourseManagement;