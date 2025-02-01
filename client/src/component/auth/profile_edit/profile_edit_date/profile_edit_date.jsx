import React, { useState, useEffect } from 'react';
import './profile_edit_date.css'; // Подключите стили, если необходимо

const ProfileEditDate = () => {
    const [userData, setUserData] = useState({
        fullName: '',
        education: '',
        residence: '',
        birthDate: '',
        homeAddress: '',
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    useEffect(() => {
        const fetchUserData = async () => {
            try {
                const response = await fetch('http://localhost:5000/profile', {
                    credentials: 'include',
                });
                if (response.ok) {
                    const data = await response.json();
                    setUserData({
                        ...userData,
                        fullName: data.fullName || '',
                        education: data.education || '',
                        residence: data.residence || '',
                        birthDate: data.birthDate || '',
                        homeAddress: data.homeAddress || '',
                    });
                } else {
                    setError('Ошибка при загрузке данных пользователя');
                }
            } catch (err) {
                setError('Ошибка при загрузке данных пользователя');
            }
        };

        fetchUserData();
    }, []);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setUserData({ ...userData, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (userData.newPassword && userData.newPassword !== userData.confirmPassword) {
            setError('Новый пароль и подтверждение пароля не совпадают');
            return;
        }

        try {
            const response = await fetch('http://localhost:5000/update-profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify(userData),
            });

            if (response.ok) {
                setSuccess('Данные успешно обновлены');
                navigate('/profile');
            } else {
                const data = await response.json();
                setError(data.message || 'Ошибка при обновлении данных');
            }
        } catch (err) {
            setError('Ошибка при обновлении данных');
        }
    };


    return (
        <div className="profile-edit-container">
            {success && <div className="success-message">{success}</div>}
            <form onSubmit={handleSubmit}>
                <div className="form-group">
                    <label>Имя Фамилия Отчество:</label>
                    <input
                        type="text"
                        name="fullName"
                        value={userData.fullName}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Образование:</label>
                    <input
                        type="text"
                        name="education"
                        value={userData.education}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Место жительства:</label>
                    <input
                        type="text"
                        name="residence"
                        value={userData.residence}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Дата рождения:</label>
                    <input
                        type="date"
                        name="birthDate"
                        value={userData.birthDate}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Домашний адрес:</label>
                    <input
                        type="text"
                        name="homeAddress"
                        value={userData.homeAddress}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Старый пароль (для изменения пароля):</label>
                    <input
                        type="password"
                        name="oldPassword"
                        value={userData.oldPassword}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Новый пароль:</label>
                    <input
                        type="password"
                        name="newPassword"
                        value={userData.newPassword}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label>Подтвердите новый пароль:</label>
                    <input
                        type="password"
                        name="confirmPassword"
                        value={userData.confirmPassword}
                        onChange={handleChange}
                    />
                </div>
                <button type="submit" className="submit-button">Сохранить изменения</button>
            </form>
        </div>
    );
};

export default ProfileEditDate;