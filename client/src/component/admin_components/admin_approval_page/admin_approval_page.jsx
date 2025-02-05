import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './admin_approval_page.css';

const AdminApprovalPage = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const [rejectReason, setRejectReason] = useState(''); // Состояние для причины отклонения
    const [selectedRegistrationId, setSelectedRegistrationId] = useState(null); // ID выбранной заявки

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
        // Устанавливаем ID выбранной заявки и открываем модальное окно
        setSelectedRegistrationId(registrationId);
    };

    const confirmReject = () => {
        if (!rejectReason) {
            setError('Укажите причину отклонения');
            return;
        }

        fetch(`http://localhost:5000/admin/reject-registration/${selectedRegistrationId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify({ reason: rejectReason }), // Отправляем причину отклонения
        })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    // Обновляем список заявок после отклонения
                    setRegistrations(registrations.map(reg =>
                        reg._id === selectedRegistrationId ? { ...reg, status: 'Отклоненный', rejectReason: rejectReason } : reg
                    ));
                    setRejectReason(''); // Очищаем поле причины
                    setSelectedRegistrationId(null); // Закрываем модальное окно
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
                            <td>{registration.courseTitle}</td>
                            <td>{registration.userName}</td>
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

            {/* Модальное окно для ввода причины отклонения */}
            {selectedRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Укажите причину отклонения</h2>
                        <textarea
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения"
                        />
                        <button onClick={confirmReject}>Подтвердить</button>
                        <button onClick={() => setSelectedRegistrationId(null)}>Отмена</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminApprovalPage;