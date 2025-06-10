import React, { useState, useEffect } from 'react';
import { Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import './admin_diagram.css';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const AdminDiagram = () => {
    const [courses, setCourses] = useState([]);
    const [registrations, setRegistrations] = useState([]);
    const [prices, setPrices] = useState([]);
    const [educations, setEducations] = useState([]);
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');
    const [triggerFetch, setTriggerFetch] = useState(true);
    const [dateError, setDateError] = useState(null);
    const [isFiltering, setIsFiltering] = useState(false);
    const [noData, setNoData] = useState(false);
    const [profitViewMode, setProfitViewMode] = useState('both');

    useEffect(() => {
        const fetchData = async () => {
            if (!triggerFetch) return;

            if (startDate && endDate && new Date(endDate) < new Date(startDate)) {
                setDateError('Конечная дата не может быть раньше начальной');
                setTriggerFetch(false);
                setIsFiltering(false);
                return;
            }

            try {
                setLoading(true);
                setIsFiltering(true);
                setDateError(null);
                setNoData(false);

                const query = startDate && endDate
                    ? `?startDate=${startDate}&endDate=${endDate}`
                    : '';

                const [regResponse, courseResponse, priceResponse, educationResponse, groupResponse] = await Promise.all([
                    fetch(`http://localhost:5000/admin/course-registrations${query}`, { credentials: 'include' }),
                    fetch(`http://localhost:5000/courses${query}`, { credentials: 'include' }),
                    fetch(`http://localhost:5000/prices${query}`, { credentials: 'include' }),
                    fetch(`http://localhost:5000/admin/educations${query}`, { credentials: 'include' }),
                    fetch(`http://localhost:5000/groups${query}`, { credentials: 'include' }),
                ]);

                if (!regResponse.ok || !courseResponse.ok || !priceResponse.ok || !educationResponse.ok || !groupResponse.ok) {
                    throw new Error(`Ошибка при загрузке данных: ${regResponse.status}, ${courseResponse.status}, ${priceResponse.status}, ${educationResponse.status}, ${groupResponse.status}`);
                }

                const [regData, courseData, priceData, educationData, groupData] = await Promise.all([
                    regResponse.json(),
                    courseResponse.json(),
                    priceResponse.json(),
                    educationResponse.json(),
                    groupResponse.json(),
                ]);

                const filteredRegistrations = startDate && endDate
                    ? regData.filter(reg => {
                        const dateField = reg.registerDate || null;
                        if (!dateField) return false;
                        const regDate = new Date(dateField);
                        if (isNaN(regDate)) return false;
                        return regDate >= new Date(startDate) && regDate <= new Date(endDate);
                    })
                    : regData;

                if (filteredRegistrations.length === 0 && startDate && endDate) {
                    setNoData(true);
                }

                setCourses(courseData || []);
                setRegistrations(filteredRegistrations || []);
                setPrices(priceData || []);
                setEducations(educationData || []);
                setGroups(groupData || []);
            } catch (err) {
                setError(err.message);
                console.error('Ошибка запроса:', err);
            } finally {
                setLoading(false);
                setTriggerFetch(false);
                setIsFiltering(false);
            }
        };

        fetchData();
    }, [triggerFetch]);

    const handleDateChange = (e) => {
        const { name, value } = e.target;
        if (name === 'startDate') {
            setStartDate(value);
        } else if (name === 'endDate') {
            setEndDate(value);
        }
        setDateError(null);
        setNoData(false);
    };

    const handleFilterClick = () => {
        setTriggerFetch(true);
    };

    const handleResetClick = () => {
        setStartDate('');
        setEndDate('');
        setDateError(null);
        setNoData(false);
        setTriggerFetch(true);
    };

    const handleLast7Days = () => {
        const end = new Date();
        const start = new Date();
        start.setDate(end.getDate() - 7);
        setStartDate(start.toISOString().split('T')[0]);
        setEndDate(end.toISOString().split('T')[0]);
        setTriggerFetch(true);
    };

    const handleViewModeChange = (mode) => {
        setProfitViewMode(mode);
    };

    const generateGradientColors = (count) => {
        const colors = [];
        const baseHue = Math.floor(Math.random() * 360);
        for (let i = 0; i < count; i++) {
            const hue = (baseHue + (i * 30)) % 360;
            colors.push(`hsl(${hue}, 70%, 60%)`);
        }
        return colors;
    };

    if (loading) return (
        <div className="loading-container">
            <div className="loading-spinner"></div>
            <span className="loading-text">Загрузка данных...</span>
        </div>
    );

    if (error) return (
        <div className="error-container">
            <span className="error-icon">⚠️</span>
            <span className="error-text">Ошибка: {error}</span>
        </div>
    );

    if (dateError) return (
        <div className="error-container">
            <span className="error-icon">❌</span>
            <span className="error-text">Ошибка: {dateError}</span>
        </div>
    );

    if (noData) return (
        <div className="no-data-container">
            <span className="no-data-icon">📉</span>
            <span className="no-data-text">Нет данных для выбранного диапазона дат</span>
            <button onClick={handleResetClick} className="reset-button">Сбросить фильтры</button>
        </div>
    );

    const statusCounts = registrations.reduce((acc, reg) => {
        const status = reg.status || 'Неизвестно';
        acc[status] = (acc[status] || 0) + 1;
        return acc;
    }, {});

    const courseTypeCounts = courses.reduce((acc, course) => {
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});

    const profitByCourse = registrations.reduce((acc, reg) => {
        if (reg.status === 'Завершил') {
            const courseTitle = reg.courseTitle || 'Неизвестный курс';
            const price = reg.price || 0;
            acc[courseTitle] = (acc[courseTitle] || 0) + Number(price);
        }
        return acc;
    }, {});

    const studentsByCourseType = registrations.reduce((acc, reg) => {
        const course = courses.find(c => c.title === reg.courseTitle) || { type: 'Без типа' };
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});

    const groupCounts = registrations.reduce((acc, reg) => {
        const groupName = reg.groupName || 'Без группы';
        acc[groupName] = (acc[groupName] || 0) + 1;
        return acc;
    }, {});

    const courseNames = Object.keys(profitByCourse);
    const profitValues = Object.values(profitByCourse);
    const backgroundColors = generateGradientColors(courseNames.length);
    const borderColors = backgroundColors.map(color => color.replace('60%)', '40%)'));

    const profitChartData = {
        labels: courseNames,
        datasets: [{
            label: 'Прибыль (руб.)',
            data: profitValues,
            backgroundColor: backgroundColors,
            borderColor: borderColors,
            borderWidth: 2,
            borderRadius: 6,
            hoverBackgroundColor: backgroundColors.map(color => color.replace('60%)', '50%)')),
            hoverBorderColor: borderColors.map(color => color.replace('40%)', '30%)')),
            hoverBorderWidth: 3
        }]
    };

    const chartOptions = {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: false
            },
            title: {
                display: true,
                text: 'Прибыль по курсам',
                font: {
                    size: 18,
                    weight: 'bold'
                },
                padding: {
                    top: 10,
                    bottom: 20
                }
            },
            tooltip: {
                backgroundColor: 'rgba(0, 0, 0, 0.8)',
                titleFont: {
                    size: 14,
                    weight: 'bold'
                },
                bodyFont: {
                    size: 12
                },
                padding: 12,
                cornerRadius: 8,
                displayColors: true,
                callbacks: {
                    label: function (context) {
                        return `${context.parsed.y.toLocaleString()} руб.`;
                    }
                }
            }
        },
        scales: {
            x: {
                display: false,
                grid: {
                    display: false
                }
            },
            y: {
                beginAtZero: true,
                grid: {
                    color: 'rgba(0, 0, 0, 0.05)'
                },
                ticks: {
                    callback: function (value) {
                        return value.toLocaleString() + ' руб.';
                    }
                },
                title: {
                    display: true,
                    text: 'Прибыль (руб.)',
                    font: {
                        weight: 'bold'
                    }
                }
            }
        },
        animation: {
            duration: 1500,
            easing: 'easeOutQuart'
        },
        elements: {
            bar: {
                borderSkipped: false
            }
        }
    };

    return (
        <div className="admin-diagram-container">
            <h1 className="admin-diagram-title">Аналитика курсов</h1>
            <div className={`date-filter-container ${isFiltering ? 'filtering' : ''}`}>
                <div className="date-input-group">
                    <label htmlFor="startDate" className="date-label">Начальная дата</label>
                    <input
                        type="date"
                        id="startDate"
                        name="startDate"
                        value={startDate}
                        onChange={handleDateChange}
                        className="date-input"
                        placeholder="Выберите дату"
                        max={new Date().toISOString().split('T')[0]}
                    />
                </div>
                <div className="date-input-group">
                    <label htmlFor="endDate" className="date-label">Конечная дата</label>
                    <input
                        type="date"
                        id="endDate"
                        name="endDate"
                        value={endDate}
                        onChange={handleDateChange}
                        className="date-input"
                        placeholder="Выберите дату"
                        max={new Date().toISOString().split('T')[0]}
                    />
                </div>
                <div className="date-button-group">
                    <button onClick={handleFilterClick} className="filter-button" disabled={isFiltering}>
                        {isFiltering ? (
                            <span className="filter-button-spinner"></span>
                        ) : (
                            'Применить'
                        )}
                    </button>
                    <button onClick={handleResetClick} className="reset-button" disabled={isFiltering}>
                        Сбросить
                    </button>
                    <button onClick={handleLast7Days} className="filter-button" disabled={isFiltering}>
                        Последние 7 дней
                    </button>
                </div>
            </div>
            {(!startDate || !endDate) && (
                <div className="info-container">
                    <span className="info-text">Данные отображаются без фильтра по датам</span>
                </div>
            )}
            <div className="admin-diagram-grid">
                <div className="diagram-card">
                    <h2 className="diagram-title">📊 Статусы регистраций</h2>
                    <div className="text-container">
                        {Object.keys(statusCounts).length === 0 ? (
                            <p className="no-data-text">Нет данных о статусах</p>
                        ) : (
                            <ul className="summary-list">
                                {Object.entries(statusCounts).map(([status, count]) => (
                                    <li key={status} className="summary-item">
                                        <span className="summary-label">{status}</span>
                                        <span className="summary-value">{count} регистраций</span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">📚 Типы курсов</h2>
                    <div className="text-container">
                        {Object.keys(courseTypeCounts).length === 0 ? (
                            <p className="no-data-text">Нет данных о типах курсов</p>
                        ) : (
                            <ul className="summary-list">
                                {Object.entries(courseTypeCounts).map(([type, count]) => (
                                    <li key={type} className="summary-item">
                                        <span className="summary-label">{type}</span>
                                        <span className="summary-value">{count} курсов</span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>

                <div className="diagram-card">
                    <h2 className="diagram-title">👩‍🎓 Студенты по типам курсов</h2>
                    <div className="text-container">
                        {Object.keys(studentsByCourseType).length === 0 ? (
                            <p className="no-data-text">Нет данных о студентах по типам</p>
                        ) : (
                            <ul className="summary-list">
                                {Object.entries(studentsByCourseType).map(([type, count]) => (
                                    <li key={type} className="summary-item">
                                        <span className="summary-label">{type}</span>
                                        <span className="summary-value">{count} студентов</span>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>
            <div className="diagram-card cards">
                <h2 className="diagram-title">Прибыль по курсам</h2>
                <div className="view-mode-buttons">
                    <button
                        onClick={() => handleViewModeChange('list')}
                        className={`toggle-button ${profitViewMode === 'list' ? 'active' : ''}`}
                    >
                        Только список
                    </button>
                    <button
                        onClick={() => handleViewModeChange('chart')}
                        className={`toggle-button ${profitViewMode === 'chart' ? 'active' : ''}`}
                    >
                        Только диаграмма
                    </button>
                    <button
                        onClick={() => handleViewModeChange('both')}
                        className={`toggle-button ${profitViewMode === 'both' ? 'active' : ''}`}
                    >
                        Список + Диаграмма
                    </button>
                </div>
                <div className="text-container">
                    {Object.keys(profitByCourse).length === 0 ? (
                        <p className="no-data-text">Нет данных о прибыли</p>
                    ) : (
                        <>
                            {(profitViewMode === 'list' || profitViewMode === 'both') && (
                                <ul className="summary-list">
                                    {courseNames.map((course, index) => (
                                        <li key={course} className="summary-item">
                                           
                                            <span className="summary-label">{course}</span>
                                            <span className="summary-value">{profitValues[index].toLocaleString()} руб.</span>
                                        </li>
                                    ))}
                                </ul>
                            )}
                            {(profitViewMode === 'chart' || profitViewMode === 'both') && (
                                <div className="chart-container" style={{ height: '600px', width: '100%' }}>
                                    <Bar data={profitChartData} options={chartOptions} />
                                </div>
                            )}
                        </>
                    )}
                </div>
            </div>
        </div>
    );
};

export default AdminDiagram;