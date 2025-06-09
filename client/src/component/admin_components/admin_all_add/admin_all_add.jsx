import React, { useState, useEffect } from 'react';
import './admin_all_add.css';

const AdminAllAdd = () => {
    const [educations, setEducations] = useState([]); // Инициализируем пустым массивом вместо null
    const [newEducation, setNewEducation] = useState('');
    const [editingId, setEditingId] = useState(null);
    const [editValue, setEditValue] = useState('');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // Загрузка списка образований
    useEffect(() => {
        const fetchEducations = async () => {
            try {
                setLoading(true);
                const response = await fetch('http://localhost:5000/admin/educations', {
                    credentials: 'include',
                });

                if (!response.ok) {
                    throw new Error('Ошибка при загрузке данных');
                }

                const data = await response.json();
                setEducations(data || []); // Гарантируем, что data будет массивом
                setError(null);
            } catch (err) {
                console.error('Ошибка при загрузке уровней образования:', err);
                setError(err.message);
                setEducations([]); // Устанавливаем пустой массив в случае ошибки
            } finally {
                setLoading(false);
            }
        };

        fetchEducations();
    }, []);

    const handleAddEducation = async () => {
        if (!newEducation.trim()) return;

        try {
            const response = await fetch('http://localhost:5000/admin/educations/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ name: newEducation }),
            });

            if (!response.ok) {
                throw new Error('Ошибка при добавлении');
            }

            const result = await response.json();
            setEducations(prev => [...prev, result]);
            setNewEducation('');
        } catch (err) {
            console.error('Ошибка при добавлении уровня образования:', err);
            setError(err.message);
        }
    };

    const startEditing = (id, name) => {
        setEditingId(id);
        setEditValue(name);
    };

    const handleUpdateEducation = async (id) => {
        if (!editValue.trim()) return;

        try {
            const response = await fetch(`http://localhost:5000/admin/educations/update/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({ name: editValue }),
            });

            if (!response.ok) {
                throw new Error('Ошибка при обновлении');
            }

            setEducations(prev =>
                prev.map(edu =>
                    edu._id === id ? { ...edu, name: editValue } : edu
                )
            );
            setEditingId(null);
        } catch (err) {
            console.error('Ошибка при обновлении уровня образования:', err);
            setError(err.message);
        }
    };

    const handleDeleteEducation = async (id) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот уровень образования?')) return;

        try {
            const response = await fetch(`http://localhost:5000/admin/educations/delete/${id}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || 'Ошибка при удалении');
            }

            setEducations(prev => prev.filter(edu => edu._id !== id));
        } catch (err) {
            console.error('Ошибка при удалении уровня образования:', err);
            setError(err.message);
            alert(err.message);
        }
    };

    if (loading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div className="error-message">Ошибка: {error}</div>;
    }

    return (
        <div className="admin-education-container">

            <div className="add-education-form">
                <input
                    type="text"
                    value={newEducation}
                    onChange={(e) => setNewEducation(e.target.value)}
                    placeholder="Новый уровень образования"
                />
                <button class='approve-btn' onClick={handleAddEducation}>Добавить</button>
            </div>

            <div className="educations-list">
                <h3>Список уровней образования</h3>
                {educations.length === 0 ? (
                    <p>Нет доступных уровней образования</p>
                ) : (
                    <ul>
                        {educations.map(edu => (
                            <li key={edu._id}>
                                {editingId === edu._id ? (
                                    <>
                                        <input
                                            type="text"
                                            value={editValue}
                                            onChange={(e) => setEditValue(e.target.value)}
                                        />
                                        <button className='approve-btn' onClick={() => handleUpdateEducation(edu._id)}>Сохранить</button>
                                        <button className='reject-btn' onClick={() => setEditingId(null)}>Отмена</button>
                                    </>
                                ) : (
                                    <>
                                        <span>{edu.name}</span>
                                        <button className='approve-btn' onClick={() => startEditing(edu._id, edu.name)}>Редактировать</button>
                                        <button className='reject-btn' onClick={() => handleDeleteEducation(edu._id)}>Удалить</button>
                                    </>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
};

export default AdminAllAdd;