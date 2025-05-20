import React, { useState, useEffect } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePriceId, setCoursePriceId] = useState('');
  const [courseTypeId, setCourseTypeId] = useState('');
  const [registrationStart, setRegistrationStart] = useState('');
  const [registrationEnd, setRegistrationEnd] = useState('');
  const [courseMaxStudents, setCourseMaxStudents] = useState(''); // Новое состояние для maxStudents
  const [courseTypes, setCourseTypes] = useState([]);
  const [prices, setPrices] = useState([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isTypeModalOpen, setIsTypeModalOpen] = useState(false);
  const [isPriceModalOpen, setIsPriceModalOpen] = useState(false);
  const [newTypeName, setNewTypeName] = useState('');
  const [editingTypeId, setEditingTypeId] = useState(null);
  const [editingTypeName, setEditingTypeName] = useState('');
  const [newPrice, setNewPrice] = useState({ amount: '', description: '' });
  const [selectedPrices, setSelectedPrices] = useState([]);
  const [percentIncrease, setPercentIncrease] = useState(0);
  const [filterType, setFilterType] = useState('Все');

  // Текущая дата для ограничения input
  const today = new Date().toISOString().split('T')[0];

  // Загрузка данных
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [coursesRes, pricesRes, typesRes] = await Promise.all([
          fetch('http://localhost:5000/courses'),
          fetch('http://localhost:5000/prices'),
          fetch('http://localhost:5000/admin/course-types', { credentials: 'include' })
        ]);

        const coursesData = await coursesRes.json();
        const pricesData = await pricesRes.json();
        const typesData = await typesRes.json();

        setCourses(coursesData || []);
        setPrices(pricesData || []);
        setCourseTypes(typesData || []);
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
        setCourses([]);
        setPrices([]);
        setCourseTypes([]);
      }
    };

    fetchData();
  }, []);

  // Обработчик Escape
  useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.key === 'Escape') {
        if (isModalOpen) setIsModalOpen(false);
        if (isPriceModalOpen) setIsPriceModalOpen(false);
        if (isTypeModalOpen) setIsTypeModalOpen(false);
        if (selectedCourse) handleCloseUpdateForm();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isModalOpen, isPriceModalOpen, isTypeModalOpen, selectedCourse]);

  // Добавление курса
  const handleAddCourse = async (e) => {
    e.preventDefault();
    if (!coursePriceId || !courseTypeId || !registrationStart || !registrationEnd || !courseMaxStudents) {
      alert('Пожалуйста, заполните все обязательные поля, включая максимальное количество студентов');
      return;
    }

    const registrationStartDate = new Date(registrationStart);
    const registrationEndDate = new Date(registrationEnd);
    const todayDate = new Date();
    todayDate.setHours(0, 0, 0, 0); // Сбрасываем время до полуночи
    const maxStudents = parseInt(courseMaxStudents, 10);

    if (isNaN(registrationStartDate.getTime()) || isNaN(registrationEndDate.getTime())) {
      alert('Некорректная дата начала или окончания регистрации');
      return;
    }

    if (registrationStartDate < todayDate) {
      alert('Дата начала регистрации не может быть раньше сегодняшней даты');
      return;
    }

    if (registrationEndDate < registrationStartDate) {
      alert('Дата окончания регистрации не может быть раньше даты начала');
      return;
    }

    if (isNaN(maxStudents) || maxStudents < 1) {
      alert('Максимальное количество студентов должно быть положительным числом');
      return;
    }

    const selectedPrice = prices.find(price => price._id === coursePriceId);
    const selectedType = courseTypes.find(type => type._id === courseTypeId);

    if (!selectedPrice || !selectedType) {
      alert('Выбранная стоимость или тип курса не найдены');
      return;
    }

    const courseData = {
      title: courseTitle,
      description: courseDescription,
      duration: parseInt(courseDuration, 10),
      priceId: selectedPrice._id,
      price: selectedPrice.amount,
      typeId: selectedType._id,
      type: selectedType.name,
      registrationStart: registrationStartDate.toISOString(),
      registrationEnd: registrationEndDate.toISOString(),
      studentsCount: 0,
      maxStudents: maxStudents, // Добавляем maxStudents
    };

    try {
      const response = await fetch('http://localhost:5000/add-course', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(courseData),
      });

      let errorMessage = 'Ошибка при добавлении курса';
      if (!response.ok) {
        try {
          const errorData = await response.json();
          errorMessage = errorData.error || errorMessage;
        } catch {
          errorMessage = await response.text() || errorMessage;
        }
        throw new Error(errorMessage);
      }

      const data = await response.json();
      alert('Курс успешно добавлен');
      setIsModalOpen(false);
      resetCourseForm();
      fetchCourses();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  // Обновление курса
  const handleUpdateCourse = async (e) => {
    e.preventDefault();
    if (!selectedCourse) return;

    if (!coursePriceId || !courseTypeId || !registrationStart || !registrationEnd || !courseMaxStudents) {
      alert('Пожалуйста, заполните все обязательные поля, включая максимальное количество студентов');
      return;
    }

    const registrationStartDate = new Date(registrationStart);
    const registrationEndDate = new Date(registrationEnd);
    const todayDate = new Date();
    todayDate.setHours(0, 0, 0, 0);
    const maxStudents = parseInt(courseMaxStudents, 10);

    if (isNaN(registrationStartDate.getTime()) || isNaN(registrationEndDate.getTime())) {
      alert('Некорректная дата начала или окончания регистрации');
      return;
    }

    if (registrationStartDate < todayDate) {
      alert('Дата начала регистрации не может быть раньше сегодняшней даты');
      return;
    }

    if (registrationEndDate < registrationStartDate) {
      alert('Дата окончания регистрации не может быть раньше даты начала');
      return;
    }

    if (isNaN(maxStudents) || maxStudents < 1) {
      alert('Максимальное количество студентов должно быть положительным числом');
      return;
    }

    if (maxStudents < selectedCourse.studentsCount) {
      alert('Максимальное количество студентов не может быть меньше текущего количества зарегистрированных');
      return;
    }

    const selectedPrice = prices.find(price => price._id === coursePriceId);
    const selectedType = courseTypes.find(type => type._id === courseTypeId);

    if (!selectedPrice || !selectedType) {
      alert('Выбранная стоимость или тип курса не найдены');
      return;
    }

    const courseData = {
      title: courseTitle,
      description: courseDescription,
      duration: parseInt(courseDuration, 10),
      priceId: selectedPrice._id,
      price: selectedPrice.amount,
      typeId: selectedType._id,
      type: selectedType.name,
      registrationStart: registrationStartDate.toISOString(),
      registrationEnd: registrationEndDate.toISOString(),
      studentsCount: selectedCourse.studentsCount,
      maxStudents: maxStudents, // Добавляем maxStudents
    };

    try {
      const response = await fetch(`http://localhost:5000/update-course/${selectedCourse._id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(courseData),
      });

      let errorMessage = 'Ошибка при обновлении курса';
      if (!response.ok) {
        try {
          const errorData = await response.json();
          errorMessage = errorData.error || errorMessage;
        } catch {
          errorMessage = await response.text() || errorMessage;
        }
        throw new Error(errorMessage);
      }

      alert('Курс успешно обновлен');
      handleCloseUpdateForm();
      fetchCourses();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  // Удаление курса
  const handleDeleteCourse = async (courseId) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот курс?')) return;

    try {
      const response = await fetch(`http://localhost:5000/delete-course/${courseId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Ошибка при удалении курса');
      }

      alert('Курс успешно удален');
      fetchCourses();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  // Работа с типами курсов
  const handleAddCourseType = async () => {
    if (!newTypeName.trim()) return;

    try {
      const response = await fetch('http://localhost:5000/admin/course-types/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: newTypeName }),
      });

      if (!response.ok) {
        throw new Error('Ошибка при добавлении типа курса');
      }

      const result = await response.json();
      setCourseTypes(prev => [...prev, { _id: result.id, name: newTypeName }]);
      setNewTypeName('');
      alert('Тип курса успешно добавлен');
    } catch (err) {
      console.error('Ошибка:', err);
      alert(err.message);
    }
  };

  const handleUpdateCourseType = async () => {
    if (!editingTypeId || !editingTypeName.trim()) return;

    try {
      const response = await fetch(`http://localhost:5000/admin/course-types/update/${editingTypeId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: editingTypeName }),
      });

      if (!response.ok) {
        throw new Error('Ошибка при обновлении типа курса');
      }

      setCourseTypes(prev =>
        prev.map(type =>
          type._id === editingTypeId ? { ...type, name: editingTypeName } : type
        )
      );
      setEditingTypeId(null);
      setEditingTypeName('');
      alert('Тип курса успешно обновлен');
    } catch (err) {
      console.error('Ошибка:', err);
      alert(err.message);
    }
  };

  const handleDeleteCourseType = async (typeId) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот тип курса?')) return;

    try {
      const response = await fetch(`http://localhost:5000/admin/course-types/delete/${typeId}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Ошибка при удалении типа курса');
      }

      setCourseTypes(prev => prev.filter(type => type._id !== typeId));
      alert('Тип курса успешно удален');
    } catch (err) {
      console.error('Ошибка:', err);
      alert(err.message);
    }
  };

  // Работа с ценами
  const handleAddPrice = async (e) => {
    e.preventDefault();
    if (!newPrice.amount || isNaN(newPrice.amount)) {
      alert('Пожалуйста, введите корректную сумму');
      return;
    }

    try {
      const response = await fetch('http://localhost:5000/add-price', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          amount: Number(newPrice.amount),
          description: newPrice.description || "",
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Ошибка при добавлении стоимости');
      }

      alert('Стоимость успешно добавлена');
      setIsPriceModalOpen(false);
      setNewPrice({ amount: '', description: '' });
      fetchPrices();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  const handleBulkUpdate = async () => {
    if (selectedPrices.length === 0 || !percentIncrease) {
      alert('Выберите стоимости и укажите процент');
      return;
    }

    if (percentIncrease <= 0) {
      alert('Процент должен быть положительным');
      return;
    }

    try {
      const response = await fetch('http://localhost:5000/bulk-update-prices', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          priceIds: selectedPrices,
          percent: parseFloat(percentIncrease),
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || 'Ошибка при массовом обновлении');
      }

      alert(`Успешно обновлено ${data.message}`);
      setSelectedPrices([]);
      setPercentIncrease(0);
      fetchPrices();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  // Вспомогательные функции
  const fetchCourses = async () => {
    try {
      const response = await fetch('http://localhost:5000/courses');
      const data = await response.json();
      setCourses(data || []);
    } catch (error) {
      console.error('Ошибка при загрузке курсов:', error);
      setCourses([]);
    }
  };

  const fetchPrices = async () => {
    try {
      const response = await fetch('http://localhost:5000/prices');
      const data = await response.json();
      setPrices(data || []);
    } catch (error) {
      console.error('Ошибка при загрузке цен:', error);
      setPrices([]);
    }
  };

  const resetCourseForm = () => {
    setCourseTitle('');
    setCourseDescription('');
    setCourseDuration('');
    setCoursePriceId('');
    setCourseTypeId('');
    setRegistrationStart('');
    setRegistrationEnd('');
    setCourseMaxStudents(''); // Сбрасываем maxStudents
    setSelectedCourse(null);
  };

  const handleCourseSelect = (course) => {
    if (!course || !course._id) {
      console.error("Ошибка: отсутствует ID курса");
      return;
    }
    setSelectedCourse(course);
    setCourseTitle(course.title);
    setCourseDescription(course.description);
    setCourseDuration(course.duration);
    setCoursePriceId(course.priceId || '');
    setCourseTypeId(course.typeId || '');
    setRegistrationStart(course.registrationStart ? new Date(course.registrationStart).toISOString().split('T')[0] : '');
    setRegistrationEnd(course.registrationEnd ? new Date(course.registrationEnd).toISOString().split('T')[0] : '');
    setCourseMaxStudents(course.maxStudents || ''); // Устанавливаем maxStudents
  };

  const handleCloseUpdateForm = () => {
    setSelectedCourse(null);
    resetCourseForm();
  };

  const togglePriceSelection = (priceId) => {
    setSelectedPrices(prev =>
      prev.includes(priceId)
        ? prev.filter(id => id !== priceId)
        : [...prev, priceId]
    );
  };

  const startEditingType = (type) => {
    setEditingTypeId(type._id);
    setEditingTypeName(type.name);
  };

  const cancelEditingType = () => {
    setEditingTypeId(null);
    setEditingTypeName('');
  };

  const filteredCourses = filterType === 'Все'
    ? courses
    : courses.filter(course => course.type === filterType);

  return (
    <div className="admin-course-management">
      <div className="admin-course-header">
        <div className="filter-section">
          <label>Фильтр по типу курса:</label>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
          >
            <option value="Все">Все</option>
            {courseTypes.map(type => (
              <option key={type._id} value={type.name}>{type.name}</option>
            ))}
          </select>
        </div>

        <div className="button-group">
          <button className='approve-btn' onClick={() => setIsModalOpen(true)}>
            Добавить курс
          </button>
          <button className='approve-btn' onClick={() => setIsPriceModalOpen(true)}>
            Управление стоимостями
          </button>
          <button className='approve-btn' onClick={() => setIsTypeModalOpen(true)}>
            Управление типами
          </button>
        </div>
      </div>

      {/* Модальное окно добавления курса */}
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
                required
              />
              <textarea
                placeholder="Описание"
                value={courseDescription}
                onChange={(e) => setCourseDescription(e.target.value)}
                required
              />
              <input
                type="number"
                placeholder="Продолжительность (часы)"
                value={courseDuration}
                onChange={(e) => setCourseDuration(e.target.value)}
                required
              />
              <input
                type="date"
                placeholder="Дата начала регистрации"
                value={registrationStart}
                onChange={(e) => setRegistrationStart(e.target.value)}
                min={today}
                required
              />
              <input
                type="date"
                placeholder="Дата окончания регистрации"
                value={registrationEnd}
                onChange={(e) => setRegistrationEnd(e.target.value)}
                min={registrationStart || today}
                required
              />
              <input
                type="number"
                placeholder="Максимальное кол-во студентов"
                value={courseMaxStudents}
                onChange={(e) => setCourseMaxStudents(e.target.value)}
                min="1"
                required
              />
              <div className="price-selection">
                <label>Выберите стоимость:</label>
                <select
                  value={coursePriceId}
                  onChange={(e) => setCoursePriceId(e.target.value)}
                  required
                >
                  <option value="">-- Выберите стоимость --</option>
                  {prices.map(price => (
                    <option key={price._id} value={price._id}>
                      {price.amount} руб. - {price.description || 'Без описания'}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label>Тип курса:</label>
                <select
                  value={courseTypeId}
                  onChange={(e) => setCourseTypeId(e.target.value)}
                  required
                >
                  <option value="">-- Выберите тип курса --</option>
                  {courseTypes.map(type => (
                    <option key={type._id} value={type._id}>
                      {type.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-buttons">
                <button className='approve-btn' type="submit">
                  Добавить
                </button>
                <button
                  className='reject-btn'
                  type="button"
                  onClick={() => {
                    setIsModalOpen(false);
                    resetCourseForm();
                  }}
                >
                  Закрыть
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Модальное окно управления стоимостями */}
      {isPriceModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Управление стоимостями</h2>

            <form onSubmit={handleAddPrice}>
              <input
                type="number"
                placeholder="Сумма (руб.)"
                value={newPrice.amount}
                onChange={(e) => setNewPrice({ ...newPrice, amount: e.target.value })}
                required
              />
              <textarea
                placeholder="Описание"
                value={newPrice.description}
                onChange={(e) => setNewPrice({ ...newPrice, description: e.target.value })}
              />
              <div className="form-buttons">
                <button className='approve-btn' type="submit">
                  Добавить стоимость
                </button>
                <button
                  className='reject-btn'
                  type="button"
                  onClick={() => setIsPriceModalOpen(false)}
                >
                  Закрыть
                </button>
              </div>
            </form>

            <div className="bulk-update-section">
              <h3>Массовое обновление</h3>
              <div className="bulk-controls">
                <input
                  type="number"
                  placeholder="Процент изменения"
                  value={percentIncrease}
                  onChange={(e) => setPercentIncrease(e.target.value)}
                />
                <button
                  className='approve-btn'
                  onClick={handleBulkUpdate}
                  disabled={selectedPrices.length === 0 || !percentIncrease}
                >
                  Применить к выбранным ({selectedPrices.length})
                </button>
              </div>
            </div>

            <div className="price-list">
              <h3>История стоимостей</h3>
              {prices.length > 0 ? (
                <table className="price-table">
                  <thead>
                    <tr>
                      <th>Выбор</th>
                      <th>Дата</th>
                      <th>Сумма</th>
                      <th>Описание</th>
                      <th>Действия</th>
                    </tr>
                  </thead>
                  <tbody>
                    {prices.map(price => (
                      <tr key={price._id}>
                        <td>
                          <input
                            type="checkbox"
                            checked={selectedPrices.includes(price._id)}
                            onChange={() => togglePriceSelection(price._id)}
                          />
                        </td>
                        <td>{new Date(price.createdAt).toLocaleDateString()}</td>
                        <td>
                          <EditableField
                            value={price.amount}
                            onSave={(newValue) => handleUpdatePrice(
                              price._id,
                              parseInt(newValue),
                              price.description
                            )}
                          />
                        </td>
                        <td>
                          <EditableField
                            value={price.description}
                            onSave={(newValue) => handleUpdatePrice(
                              price._id,
                              price.amount,
                              newValue
                            )}
                          />
                        </td>
                        <td>
                          <button
                            className='reject-btn'
                            onClick={() => handleDeletePrice(price._id)}
                          >
                            Удалить
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              ) : (
                <p>Нет сохраненных стоимостей</p>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Модальное окно управления типами курсов */}
      {isTypeModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Управление типами курсов</h2>

            <form onSubmit={(e) => {
              e.preventDefault();
              handleAddCourseType();
            }}>
              <input
                type="text"
                value={newTypeName}
                onChange={(e) => setNewTypeName(e.target.value)}
                placeholder="Новый тип курса"
                required
              />
              <div className="form-buttons">
                <button className='approve-btn' type="submit">
                  Добавить тип
                </button>
              </div>
            </form>

            <div className="types-table-container">
              <h3>Список типов курсов</h3>
              {courseTypes.length === 0 ? (
                <p>Нет доступных типов курсов</p>
              ) : (
                <table className="types-table">
                  <thead>
                    <tr>
                      <th>Название</th>
                      <th>Действия</th>
                    </tr>
                  </thead>
                  <tbody>
                    {courseTypes.map(type => (
                      <tr key={type._id}>
                        <td>
                          {editingTypeId === type._id ? (
                            <input
                              type="text"
                              value={editingTypeName}
                              onChange={(e) => setEditingTypeName(e.target.value)}
                              autoFocus
                            />
                          ) : (
                            <span>{type.name}</span>
                          )}
                        </td>
                        <td>
                          {editingTypeId === type._id ? (
                            <>
                              <button
                                className='approve-btn'
                                onClick={handleUpdateCourseType}
                              >
                                Сохранить
                              </button>
                              <button
                                className='reject-btn'
                                onClick={cancelEditingType}
                              >
                                Отмена
                              </button>
                            </>
                          ) : (
                            <>
                              <button
                                className='edit-btn'
                                onClick={() => startEditingType(type)}
                              >
                                Редактировать
                              </button>
                              <button
                                className='reject-btn'
                                onClick={() => handleDeleteCourseType(type._id)}
                              >
                                Удалить
                              </button>
                            </>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>

            <div className="modal-footer">
              <button
                className='reject-btn'
                onClick={() => {
                  setIsTypeModalOpen(false);
                  cancelEditingType();
                }}
              >
                Закрыть
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Таблица курсов */}
      <table className="course-table">
        <thead>
          <tr>
            <th>Заголовок</th>
            <th>Описание</th>
            <th>Продолжительность (часы)</th>
            <th>Стоимость (руб.)</th>
            <th>Тип</th>
            <th>Начало регистрации</th>
            <th>Окончание регистрации</th>
            <th>Студентов (тек./макс.)</th>
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
                <td>{course.registrationStart ? new Date(course.registrationStart).toLocaleDateString() : '-'}</td>
                <td>{course.registrationEnd ? new Date(course.registrationEnd).toLocaleDateString() : '-'}</td>
                <td>{course.studentsCount || 0}/{course.maxStudents || '∞'}</td>
                <td>
                  <button
                    className='reject-btn'
                    onClick={(e) => {
                      e.stopPropagation();
                      handleDeleteCourse(course._id);
                    }}
                  >
                    Удалить
                  </button>
                </td>
              </tr>
              {selectedCourse && selectedCourse._id === course._id && (
                <tr>
                  <td colSpan="9">
                    <form onSubmit={handleUpdateCourse} className="update-form">
                      <input
                        type="text"
                        placeholder="Заголовок"
                        value={courseTitle}
                        onChange={(e) => setCourseTitle(e.target.value)}
                        required
                      />
                      <textarea
                        placeholder="Описание"
                        value={courseDescription}
                        onChange={(e) => setCourseDescription(e.target.value)}
                        required
                      />
                      <input
                        type="number"
                        placeholder="Продолжительность (часы)"
                        value={courseDuration}
                        onChange={(e) => setCourseDuration(e.target.value)}
                        required
                      />
                      <input
                        type="date"
                        placeholder="Дата начала регистрации"
                        value={registrationStart}
                        onChange={(e) => setRegistrationStart(e.target.value)}
                        min={today}
                        required
                      />
                      <input
                        type="date"
                        placeholder="Дата окончания регистрации"
                        value={registrationEnd}
                        onChange={(e) => setRegistrationEnd(e.target.value)}
                        min={registrationStart || today}
                        required
                      />
                      <input
                        type="number"
                        placeholder="Максимальное кол-во студентов"
                        value={courseMaxStudents}
                        onChange={(e) => setCourseMaxStudents(e.target.value)}
                        min={course.studentsCount || 1}
                        required
                      />
                      <div className="price-selection">
                        <label>Выберите стоимость:</label>
                        <select
                          value={coursePriceId}
                          onChange={(e) => setCoursePriceId(e.target.value)}
                          required
                        >
                          <option value="">-- Выберите стоимость --</option>
                          {prices.map(price => (
                            <option key={price._id} value={price._id}>
                              {price.amount} руб. - {price.description || 'Без описания'}
                            </option>
                          ))}
                        </select>
                      </div>

                      <div className="form-group">
                        <label>Тип курса:</label>
                        <select
                          value={courseTypeId}
                          onChange={(e) => setCourseTypeId(e.target.value)}
                          required
                        >
                          <option value="">-- Выберите тип курса --</option>
                          {courseTypes.map(type => (
                            <option key={type._id} value={type._id}>
                              {type.name}
                            </option>
                          ))}
                        </select>
                      </div>

                      <div className="form-buttons">
                        <button className='approve-btn' type="submit">
                          Обновить курс
                        </button>
                        <button
                          className='reject-btn'
                          type="button"
                          onClick={handleCloseUpdateForm}
                        >
                          Закрыть
                        </button>
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

// Компонент для редактируемых полей
const EditableField = ({ value, onSave }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [currentValue, setCurrentValue] = useState(value);

  const handleSave = () => {
    onSave(currentValue);
    setIsEditing(false);
  };

  return isEditing ? (
    <div className="editable-field">
      <input
        type={typeof value === 'number' ? 'number' : 'text'}
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        autoFocus
      />
      <button className='approve-btn' onClick={handleSave}>✓</button>
      <button className='reject-btn' onClick={() => {
        setCurrentValue(value);
        setIsEditing(false);
      }}>✗</button>
    </div>
  ) : (
    <div onClick={() => setIsEditing(true)} className="editable-field">
      {value}
    </div>
  );
};

export default AdminCourseManagement;