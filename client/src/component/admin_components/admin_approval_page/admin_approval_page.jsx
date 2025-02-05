import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom'; // Используем useNavigate вместо useHistory
import './admin_approval_page.css';

const AdminApprovalPage = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const navigate = useNavigate(); // Используем useNavigate

    useEffect(() => {
        fetch('http://localhost:5000/admin/course-registrations', {
            credentials: 'include'
        })
            .then(response => response.json())
            .then(data => {
                console.log("Data from backend:", data); // Проверьте структуру данных
                if (data.error) {
                    setError(data.error);
                } else {
                    setRegistrations(data);

                }
                setIsLoading(false);
            })
            .catch(error => {
                console.error('Error fetching registrations:', error);
                setError('Ошибка при загрузке заявок');
                setIsLoading(false);
            });
    }, []);

    const handleApprove = (registrationId) => {
        fetch(`http://localhost:5000/admin/approve-registration/${registrationId}`, {
            method: 'POST',
            credentials: 'include'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Обновляем список заявок после одобрения
                    setRegistrations(registrations.map(reg =>
                        reg._id === registrationId ? { ...reg, status: 'Одобренный' } : reg
                    ));
                    console.log(registrations.map)
                } else {
                    setError(data.message || 'Ошибка при одобрении заявки');
                }
            })
            .catch(error => {
                console.error('Error approving registration:', error);
                setError('Ошибка при одобрении заявки');
            });
    };

    const handleReject = (registrationId) => {
        fetch(`http://localhost:5000/admin/reject-registration/${registrationId}`, {
            method: 'POST',
            credentials: 'include'
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Обновляем список заявок после отклонения
                    setRegistrations(registrations.map(reg =>
                        reg._id === registrationId ? { ...reg, status: 'Отклоненный' } : reg
                    ));
                    console.log(registrations.map)
                } else {
                    setError(data.message || 'Ошибка при отклонении заявки');
                }
            })
            .catch(error => {
                console.error('Error rejecting registration:', error);
                setError('Ошибка при отклонении заявки');
            });
    };

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }


    return (
        <div className="admin-approval-page">
            <h1>Заявки на курсы</h1>
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {registrations.map(registration => (
                        <tr key={registration._id}>
                            <td>{registration.courseTitle}</td> {/* Используем courseTitle */}
                            <td>{registration.userName}</td>   {/* Используем userName */}
                            <td>{registration.status}</td>
                            <td>
                                {registration.status === 'Ожидание' && (
                                    <>
                                        <button onClick={() => handleApprove(registration._id)}>Одобрить</button>
                                        <button onClick={() => handleReject(registration._id)}>Отклонить</button>
                                    </>
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default AdminApprovalPage;