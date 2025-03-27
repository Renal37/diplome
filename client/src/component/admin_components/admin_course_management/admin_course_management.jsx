import React, { useState, useEffect } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  // Состояния для курсов
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePriceId, setCoursePriceId] = useState('');
  const [courseType, setCourseType] = useState('Повышение квалификации');
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Состояния для стоимостей
  const [isPriceModalOpen, setIsPriceModalOpen] = useState(false);
  const [prices, setPrices] = useState([]);
  const [newPrice, setNewPrice] = useState({
    amount: '',
    description: ''
  });
  const [selectedPrices, setSelectedPrices] = useState([]);
  const [percentIncrease, setPercentIncrease] = useState(0);

  // Фильтрация
  const [filterType, setFilterType] = useState('Все');

  // Загрузка данных
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

  useEffect(() => {
    fetchCourses();
    fetchPrices();
  }, []);

  // Обработчики для курсов
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
    setCourseType(course.type);
  };

  const handleAddCourse = async (e) => {
    e.preventDefault();
    if (!coursePriceId) {
      alert('Пожалуйста, выберите стоимость');
      return;
    }

    const selectedPrice = prices.find(price => price._id === coursePriceId);
    if (!selectedPrice) {
      alert('Выбранная стоимость не найдена');
      return;
    }

    try {
      const response = await fetch('http://localhost:5000/add-course', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: courseTitle,
          description: courseDescription,
          duration: parseInt(courseDuration, 10),
          priceId: coursePriceId,
          price: selectedPrice.amount,
          type: courseType,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Ошибка при добавлении курса');
      }

      alert('Курс успешно добавлен');
      setIsModalOpen(false);
      resetCourseForm();
      fetchCourses();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  const handleDeleteCourse = async (courseId) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот курс?')) {
      return;
    }

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

  const resetCourseForm = () => {
    setCourseTitle('');
    setCourseDescription('');
    setCourseDuration('');
    setCoursePriceId('');
    setCourseType('Повышение квалификации');
    setSelectedCourse(null);
  };

  // Обработчики для стоимостей
  const handleAddPrice = async (e) => {
    e.preventDefault();
    if (!newPrice.amount || isNaN(newPrice.amount)) {
      alert('Пожалуйста, введите корректную сумму');
      return;
    }

    try {
      const response = await fetch('http://localhost:5000/add-price', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
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

  const handleUpdatePrice = async (priceId, newAmount, newDescription) => {
    try {
      const response = await fetch(`http://localhost:5000/update-price/${priceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          amount: newAmount,
          description: newDescription,
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || 'Ошибка при обновлении стоимости');
      }

      alert('Стоимость успешно обновлена');
      fetchPrices();
    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message);
    }
  };

  const handleDeletePrice = async (priceId) => {
    if (!window.confirm('Вы уверены, что хотите удалить эту стоимость?')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:5000/delete-price/${priceId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Ошибка при удалении стоимости');
      }

      alert('Стоимость успешно удалена');
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
        headers: {
          'Content-Type': 'application/json',
        },
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


  const togglePriceSelection = (priceId) => {
    setSelectedPrices(prev =>
      prev.includes(priceId)
        ? prev.filter(id => id !== priceId)
        : [...prev, priceId]
    );
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
            <option value="Повышение квалификации">Повышение квалификации</option>
            <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
          </select>
        </div>

        <div className="button-group">
          <button className='approve-btn' onClick={() => setIsModalOpen(true)}>
            Добавить курс
          </button>
          {/* <button className='approve-btn' onClick={() => setIsPriceModalOpen(true)}>
            Управление стоимостями
          </button> */}
        </div>
      </div>

      {/* Модальное окно добавления курса */}
      {isModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>{selectedCourse ? 'Редактировать курс' : 'Добавить курс'}</h2>
            <form onSubmit={selectedCourse ? handleUpdateCourse : handleAddCourse}>
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
                      {price.amount} руб. ({new Date(price.createdAt).toLocaleDateString()}) - {price.description}
                    </option>
                  ))}
                </select>
                <button
                  type="button"
                  className='approve-btn'
                  onClick={() => {
                    setIsModalOpen(false);
                    setIsPriceModalOpen(true);
                  }}
                >
                  Добавить новую стоимость
                </button>
              </div>

              <select
                value={courseType}
                onChange={(e) => setCourseType(e.target.value)}
                required
              >
                <option value="Повышение квалификации">Повышение квалификации</option>
                <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
              </select>

              <div className="form-buttons">
                <button className='approve-btn' type="submit">
                  {selectedCourse ? 'Обновить' : 'Добавить'}
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

            {/* Форма добавления новой стоимости */}
            <form onSubmit={handleAddPrice}>
              <input
                type="number"
                placeholder="Сумма (руб.)"
                value={newPrice.amount}
                onChange={(e) => setNewPrice({ ...newPrice, amount: e.target.value })}
                required
              />
              <textarea
                placeholder="Описание (например: 'Стоимость на 2023 год')"
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

            {/* Массовое обновление стоимостей */}
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

            {/* Таблица стоимостей */}
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

      {/* Таблица курсов */}
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
            <tr key={course._id} onClick={() => handleCourseSelect(course)}>
              <td>{course.title}</td>
              <td className="course-description">{course.description}</td>
              <td>{course.duration}</td>
              <td>{course.price}</td>
              <td>{course.type}</td>
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