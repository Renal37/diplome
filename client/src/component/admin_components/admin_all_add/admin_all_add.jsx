import React, { useState, useEffect } from 'react';
import './admin_all_add.css';

const AdminAllAdd = () => {
    const [educations, setEducations] = useState([]);
    const [newEducation, setNewEducation] = useState('');
    const [editingId, setEditingId] = useState(null);
    const [editValue, setEditValue] = useState('');

    // Загрузка списка образований
    useEffect(() => {
        fetchEducations();
    }, []);

    const fetchEducations = async () => {
        try {
            const response = await fetch('http://localhost:5000/admin/educations', {
                credentials: 'include',
            });
            if (response.ok) {
                const data = await response.json();
                setEducations(data);
            }
        } catch (err) {
            console.error('Ошибка при загрузке уровней образования:', err);
        }
    };

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

            if (response.ok) {
                setNewEducation('');
                fetchEducations();
            }
        } catch (err) {
            console.error('Ошибка при добавлении уровня образования:', err);
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

            if (response.ok) {
                setEditingId(null);
                fetchEducations();
            }
        } catch (err) {
            console.error('Ошибка при обновлении уровня образования:', err);
        }
    };

    const handleDeleteEducation = async (id) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот уровень образования?')) return;

        try {
            const response = await fetch(`http://localhost:5000/admin/educations/delete/${id}`, {
                method: 'DELETE',
                credentials: 'include',
            });

            if (response.ok) {
                fetchEducations();
            } else {
                const data = await response.json();
                alert(data.message || 'Ошибка при удалении');
            }
        } catch (err) {
            console.error('Ошибка при удалении уровня образования:', err);
        }
    };

    return (
        <div className="admin-education-container">
            <h2>Управление уровнями образования</h2>
            
            <div className="add-education-form">
                <input
                    type="text"
                    value={newEducation}
                    onChange={(e) => setNewEducation(e.target.value)}
                    placeholder="Новый уровень образования"
                />
                <button onClick={handleAddEducation}>Добавить</button>
            </div>

            <div className="educations-list">
                <h3>Список уровней образования</h3>
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
                                    <button onClick={() => handleUpdateEducation(edu._id)}>Сохранить</button>
                                    <button onClick={() => setEditingId(null)}>Отмена</button>
                                </>
                            ) : (
                                <>
                                    <span>{edu.name}</span>
                                    <button onClick={() => startEditing(edu._id, edu.name)}>Редактировать</button>
                                    <button onClick={() => handleDeleteEducation(edu._id)}>Удалить</button>
                                </>
                            )}
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
};

export default AdminAllAdd;