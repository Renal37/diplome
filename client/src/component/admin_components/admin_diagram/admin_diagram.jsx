import React, { useState, useEffect } from 'react';
import './admin_diagram.css';

const AdminDiagram = () => {
    const [courses, setCourses] = useState([]);
    const [registrations, setRegistrations] = useState([]);
    const [prices, setPrices] = useState([]);
    const [educations, setEducations] = useState([]);
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const [regResponse, courseResponse, priceResponse, educationResponse, groupResponse] = await Promise.all([
                    fetch('http://localhost:5000/admin/course-registrations', { credentials: 'include' }),
                    fetch('http://localhost:5000/courses', { credentials: 'include' }),
                    fetch('http://localhost:5000/prices', { credentials: 'include' }),
                    fetch('http://localhost:5000/admin/educations', { credentials: 'include' }),
                    fetch('http://localhost:5000/groups', { credentials: 'include' }),
                ]);

                if (!regResponse.ok || !courseResponse.ok || !priceResponse.ok || !educationResponse.ok || !groupResponse.ok) {
                    throw new Error('Ошибка при загрузке данных');
                }

                const [regData, courseData, priceData, educationData, groupData] = await Promise.all([
                    regResponse.json(),
                    courseResponse.json(),
                    priceResponse.json(),
                    educationResponse.json(),
                    groupResponse.json(),
                ]);

                setCourses(courseData);
                setRegistrations(regData);
                setPrices(priceData);
                setEducations(educationData);
                setGroups(groupData.groups || []);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

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

    // 1. Распределение статусов регистраций
    const statusCounts = registrations.reduce((acc, reg) => {
        const status = reg.status || 'Неизвестно';
        acc[status] = (acc[status] || 0) + 1;
        return acc;
    }, {});

    // 2. Курсы по типам
    const courseTypeCounts = courses.reduce((acc, course) => {
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});

    // 3. Прибыль по завершённым курсам
    const profitByCourse = registrations.reduce((acc, reg) => {
        if (reg.status === 'Завершил') {
            const course = courses.find(c => c.title === reg.courseTitle) || {};
            const price = prices.find(p => p._id === course.priceId) || { amount: course.price || 0 };
            const courseTitle = reg.courseTitle || 'Неизвестный курс';
            acc[courseTitle] = (acc[courseTitle] || 0) + (price.amount || 0);
        }
        return acc;
    }, {});

    // 4. Распределение студентов по типам курсов
    const studentsByCourseType = registrations.reduce((acc, reg) => {
        const course = courses.find(c => c.title === reg.courseTitle) || {};
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});

    // 5. Студенты по группам
    const groupCounts = registrations.reduce((acc, reg) => {
        const group = groups.find(g => g._id === reg.groupId) || { groupName: 'Без группы' };
        acc[group.groupName] = (acc[group.groupName] || 0) + 1;
        return acc;
    }, {});

    return (
        <div className="admin-diagram-container">
            <div className="admin-diagram-grid">
                <div className="diagram-card">
                    <h2 className="diagram-title">📊 Статусы регистраций</h2>
                    <div className="text-container">
                        <ul className="summary-list">
                            {Object.entries(statusCounts).map(([status, count]) => (
                                <li key={status} className="summary-item">
                                    <span className="summary-label">{status}</span>
                                    <span className="summary-value">{count} регистраций</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">📚 Типы курсов</h2>
                    <div className="text-container">
                        <ul className="summary-list">
                            {Object.entries(courseTypeCounts).map(([type, count]) => (
                                <li key={type} className="summary-item">
                                    <span className="summary-label">{type}</span>
                                    <span className="summary-value">{count} курсов</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">💰 Прибыль по курсам</h2>
                    <div className="text-container">
                        <ul className="summary-list">
                            {Object.entries(profitByCourse).map(([course, profit]) => (
                                <li key={course} className="summary-item">
                                    <span className="summary-label">{course}</span>
                                    <span className="summary-value">{profit.toLocaleString()} руб.</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">👩‍🎓 Студенты по типам курсов</h2>
                    <div className="text-container">
                        <ul className="summary-list">
                            {Object.entries(studentsByCourseType).map(([type, count]) => (
                                <li key={type} className="summary-item">
                                    <span className="summary-label">{type}</span>
                                    <span className="summary-value">{count} студентов</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">👥 Студенты по группам</h2>
                    <div className="text-container">
                        <ul className="summary-list">
                            {Object.entries(groupCounts).map(([group, count]) => (
                                <li key={group} className="summary-item">
                                    <span className="summary-label">{group}</span>
                                    <span className="summary-value">{count} студентов</span>
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AdminDiagram;