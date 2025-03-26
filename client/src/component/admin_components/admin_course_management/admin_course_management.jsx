import React, { useState, useEffect } from 'react';
import './admin_course_management.css';

const AdminCourseManagement = () => {
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState(null);
  const [courseTitle, setCourseTitle] = useState('');
  const [courseDescription, setCourseDescription] = useState('');
  const [courseDuration, setCourseDuration] = useState('');
  const [coursePriceId, setCoursePriceId] = useState('');
  const [courseType, setCourseType] = useState('Повышение квалификации');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isPriceModalOpen, setIsPriceModalOpen] = useState(false);
  const [filterType, setFilterType] = useState('Все');
  const [prices, setPrices] = useState([]);
  const [newPrice, setNewPrice] = useState({
    amount: '',
    description: ''
  });

  // Fetch courses and prices
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

  const handleAddPrice = async (e) => {
    e.preventDefault();
    
    // Валидация данных
    if (!newPrice.amount || isNaN(newPrice.amount)) {
      alert("Пожалуйста, введите корректную сумму");
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
  
      // Логирование для отладки
      console.log("Response status:", response.status);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || "Ошибка сервера");
      }
  
      const data = await response.json();
      alert(data.message || 'Стоимость успешно добавлена');
      setIsPriceModalOpen(false);
      setNewPrice({ amount: '', description: '' });
      fetchPrices();
    } catch (error) {
      console.error('Ошибка при добавлении стоимости:', error);
      alert(`Ошибка: ${error.message}`);
    }
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
    
    if (response.ok) {
      alert('Курс был успешно добавлен');
      setIsModalOpen(false);
      setCourseTitle('');
      setCourseDescription('');
      setCourseDuration('');
      setCoursePriceId('');
      setCourseType('Повышение квалификации');
      fetchCourses();
    } else {
      alert('Ошибка при добавлении курса');
    }
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
          <button className='approve-btn' onClick={() => setIsModalOpen(true)}>Добавить курс</button>
          <button className='approve-btn' onClick={() => setIsPriceModalOpen(true)}>Управление стоимостями</button>
        </div>
      </div>

      {/* Price Management Modal */}
      {isPriceModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Управление стоимостями</h2>
            <form onSubmit={handleAddPrice}>
              <input
                type="number"
                placeholder="Сумма (руб.)"
                value={newPrice.amount}
                onChange={(e) => setNewPrice({...newPrice, amount: e.target.value})}
                required
              />
              <textarea
                placeholder="Описание (например: 'Стоимость на 2023 год')"
                value={newPrice.description}
                onChange={(e) => setNewPrice({...newPrice, description: e.target.value})}
              />
              <div className="form-buttons">
                <button className='reject-btn' type="button" onClick={() => setIsPriceModalOpen(false)}>Закрыть</button>
                <button className='approve-btn' type="submit">Добавить стоимость</button>
              </div>
            </form>

            <div className="price-list">
              <h3>История стоимостей</h3>
              {prices.length > 0 ? (
                <table className="price-table">
                  <thead>
                    <tr>
                      <th>Дата создания</th>
                      <th>Сумма (руб.)</th>
                      <th>Описание</th>
                    </tr>
                  </thead>
                  <tbody>
                    {prices.map(price => (
                      <tr key={price._id}>
                        <td>{new Date(price.createdAt).toLocaleDateString()}</td>
                        <td>{price.amount}</td>
                        <td>{price.description}</td>
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

      {/* Course Creation Modal */}
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
              <input
                type="text"
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
                <button className='reject-btn' type="button" onClick={() => setIsModalOpen(false)}>Закрыть</button>
                <button className='approve-btn' type="submit">Добавить курс</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Course Table */}
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
                <button className='reject-btn' onClick={(e) => {
                  e.stopPropagation();
                  handleDeleteCourse(course._id);
                }}>Удалить</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminCourseManagement;