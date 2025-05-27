import React, { useState, useEffect } from 'react';
import { Chart as ChartJS, ArcElement, BarElement, CategoryScale, LinearScale, Title, Tooltip, Legend } from 'chart.js';
import { Pie, Bar } from 'react-chartjs-2';
import './admin_diagram.css';

// Регистрация компонентов Chart.js
ChartJS.register(ArcElement, BarElement, CategoryScale, LinearScale, Title, Tooltip, Legend);

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

                console.log('Данные курсов:', courseData);
                console.log('Регистрации:', regData);
                console.log('Цены:', priceData);
                console.log('Образование:', educationData);
                console.log('Группы:', groupData);

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

    if (loading) return <div className="loading-container"><span className="loading-text">Загрузка...</span></div>;
    if (error) return <div className="error-container"><span className="error-text">Ошибка: {error}</span></div>;

    // 1. Распределение статусов регистраций (Круговая диаграмма)
    const statusCounts = registrations.reduce((acc, reg) => {
        const status = reg.status || 'Неизвестно';
        acc[status] = (acc[status] || 0) + 1;
        return acc;
    }, {});
    const statusData = {
        labels: Object.keys(statusCounts),
        datasets: [{
            data: Object.values(statusCounts),
            backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'],
        }],
    };

    // 2. Курсы по типам (Столбчатая диаграмма)
    const courseTypeCounts = courses.reduce((acc, course) => {
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});
    const courseTypeData = {
        labels: Object.keys(courseTypeCounts),
        datasets: [{
            label: 'Количество курсов',
            data: Object.values(courseTypeCounts),
            backgroundColor: '#36A2EB',
        }],
    };

    // 3. Прибыль по завершённым курсам (Круговая диаграмма)
    const profitByCourse = registrations.reduce((acc, reg) => {
        if (reg.status === 'Завершил') {
            const course = courses.find(c => c.title === reg.courseTitle) || {};
            const price = prices.find(p => p._id === course.priceId) || { amount: course.price || 0 };
            const courseTitle = reg.courseTitle || 'Неизвестный курс';
            acc[courseTitle] = (acc[courseTitle] || 0) + (price.amount || 0);
        }
        return acc;
    }, {});
    const profitData = {
        labels: Object.keys(profitByCourse),
        datasets: [{
            data: Object.values(profitByCourse),
            backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'],
        }],
    };

    // 4. Распределение студентов по типам курсов (Столбчатая диаграмма)
    const studentsByCourseType = registrations.reduce((acc, reg) => {
        const course = courses.find(c => c.title === reg.courseTitle) || {};
        const type = course.type || 'Без типа';
        acc[type] = (acc[type] || 0) + 1;
        return acc;
    }, {});
    const studentsByCourseTypeData = {
        labels: Object.keys(studentsByCourseType),
        datasets: [{
            label: 'Количество студентов',
            data: Object.values(studentsByCourseType),
            backgroundColor: '#FFCE56',
        }],
    };

    // 5. Студенты по группам (Столбчатая диаграмма)
    const groupCounts = registrations.reduce((acc, reg) => {
        const group = groups.find(g => g._id === reg.groupId) || { groupName: 'Без группы' };
        acc[group.groupName] = (acc[group.groupName] || 0) + 1;
        return acc;
    }, {});
    const groupData = {
        labels: Object.keys(groupCounts),
        datasets: [{
            label: 'Количество студентов',
            data: Object.values(groupCounts),
            backgroundColor: '#4BC0C0',
        }],
    };

    return (
        <div className="admin-diagram-container">
            <h1 className="admin-diagram-title">Админ-панель</h1>
            <div className="admin-diagram-grid">
                <div className="diagram-card">
                    <h2 className="diagram-title">Распределение статусов регистраций</h2>
                    <div className="chart-container">
                        <Pie data={statusData} options={{ responsive: true, plugins: { legend: { position: 'top' }, title: { display: true, text: 'Распределение статусов регистраций' } } }} />
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">Курсы по типам</h2>
                    <div className="chart-container">
                        <Bar data={courseTypeData} options={{ responsive: true, plugins: { legend: { position: 'top' }, title: { display: true, text: 'Курсы по типам' } } }} />
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">Прибыль по завершённым курсам</h2>
                    <div className="chart-container">
                        <Pie data={profitData} options={{ responsive: true, plugins: { legend: { position: 'top' }, title: { display: true, text: 'Прибыль по завершённым курсам' } } }} />
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">Распределение студентов по типам курсов</h2>
                    <div className="chart-container">
                        <Bar data={studentsByCourseTypeData} options={{ responsive: true, plugins: { legend: { position: 'top' }, title: { display: true, text: 'Распределение студентов по типам курсов' } } }} />
                    </div>
                </div>
                <div className="diagram-card">
                    <h2 className="diagram-title">Студенты по группам</h2>
                    <div className="chart-container">
                        <Bar data={groupData} options={{ responsive: true, plugins: { legend: { position: 'top' }, title: { display: true, text: 'Студенты по группам' } } }} />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AdminDiagram;