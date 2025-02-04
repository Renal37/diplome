import React, { useState, useEffect } from 'react';
import './profile_edit_date.css';

const ProfileEditDate = () => {
    const [userData, setUserData] = useState({
        lastName: '',
        firstName: '',
        middleName: '',
        education: '',
        residence: '',
        birthDate: '',
        homeAddress: '',
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
        agreeToProcessing: false,
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [showPassword, setShowPassword] = useState(false);

    // Загрузка данных пользователя при монтировании компонента
    useEffect(() => {
        const fetchUserData = async () => {
            try {
                const response = await fetch('http://localhost:5000/profile', {
                    credentials: 'include', // Для отправки куки с токеном
                });
                if (response.ok) {
                    const data = await response.json();
                    console.log('Данные с сервера:', data); // Логируем данные для отладки

                    // Разделяем fullName на отдельные поля
                    const nameParts = data.fullName ? data.fullName.split(' ') : ['', '', ''];

                    // Обновляем состояние, сохраняя предыдущие значения
                    setUserData((prevState) => ({
                        ...prevState,
                        lastName: data.lastName || '',
                        firstName: data.firstName || '',
                        middleName: data.middleName || '',
                        education: data.education || '',
                        residence: data.residence || '',
                        birthDate: data.birthDate || '',
                        homeAddress: data.homeAddress || '',
                        agreeToProcessing: data.agreeToProcessing || false,
                    }));
                } else {
                    setError('Ошибка при загрузке данных пользователя');
                }
            } catch (err) {
                setError('Ошибка при загрузке данных пользователя');
                console.error('Ошибка при загрузке данных:', err);
            }
        };

        fetchUserData();
    }, []);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setUserData((prevState) => ({
            ...prevState,
            [name]: type === 'checkbox' ? checked : value,
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        if (!userData.agreeToProcessing) {
            setError('Необходимо согласие на обработку данных');
            return;
        }

        if (userData.newPassword && userData.newPassword !== userData.confirmPassword) {
            setError('Новый пароль и подтверждение пароля не совпадают');
            return;
        }

        // Собираем fullName из отдельных полей
        const fullName = `${userData.lastName} ${userData.firstName} ${userData.middleName}`.trim();

        // Создаем объект для отправки только измененных данных
        const updateData = {
            fullName,
            lastName: userData.lastName,
            firstName: userData.firstName,
            middleName: userData.middleName,
            education: userData.education,
            residence: userData.residence,
            birthDate: userData.birthDate,
            homeAddress: userData.homeAddress,
            agreeToProcessing: userData.agreeToProcessing,
        };

        // Добавляем пароль только если он был изменен
        if (userData.oldPassword && userData.newPassword) {
            updateData.oldPassword = userData.oldPassword;
            updateData.newPassword = userData.newPassword;
        }

        try {
            const response = await fetch('http://localhost:5000/update-profile', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify(updateData),
            });

            if (response.ok) {
                setSuccess('Данные успешно обновлены');
            } else {
                const data = await response.json();
                setError(data.message || 'Ошибка при обновлении данных');
            }
        } catch (err) {
            setError('Ошибка при обновлении данных');
        }
    };

    const toggleShowPassword = () => {
        setShowPassword(!showPassword);
    };

    return (
        <div className="profile-edit-container">
            <form onSubmit={handleSubmit}>
                <div className="from_input_group">
                    <div className="input_group">
                        <div className="form-group">
                            <label>Фамилия:</label>
                            <input
                                type="text"
                                name="lastName"
                                value={userData.lastName}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Имя:</label>
                            <input
                                type="text"
                                name="firstName"
                                value={userData.firstName}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Отчество:</label>
                            <input
                                type="text"
                                name="middleName"
                                value={userData.middleName}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Образование:</label>
                            <select
                                name="education"
                                value={userData.education}
                                onChange={handleChange}
                            >
                                <option value="">Выберите образование</option>
                                <option value="Высшее">Высшее</option>
                                <option value="Среднее профессиональное">Среднее профессиональное</option>
                                <option value="Среднее общее">Среднее общее</option>
                            </select>
                        </div>
                    </div>
                    <div className="input_group">
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
                            <label>Место жительства:</label>
                            <input
                                type="text"
                                name="residence"
                                value={userData.residence}
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
                        {/* {!userData.agreeToProcessing && ( */}
                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    name="agreeToProcessing"
                                    checked={userData.agreeToProcessing}
                                    onChange={handleChange}
                                />
                                Согласен на обработку персональных данных
                            </label>
                        </div>
                        {/* )} */}
                        {/* <div className="form-group">
                            <label>Старый пароль (для изменения пароля):</label>
                            <input
                                type={showPassword ? "text" : "password"}
                                name="oldPassword"
                                value={userData.oldPassword}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Новый пароль:</label>
                            <input
                                type={showPassword ? "text" : "password"}
                                name="newPassword"
                                value={userData.newPassword}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Подтвердите новый пароль:</label>
                            <input
                                type={showPassword ? "text" : "password"}
                                name="confirmPassword"
                                value={userData.confirmPassword}
                                onChange={handleChange}
                            />
                        </div>
                        <button type="button" onClick={toggleShowPassword}>
                            {showPassword ? "Скрыть пароль" : "Показать пароль"}
                        </button> */}
                    </div>
                </div>
                {error && <div className="error-message">{error}</div>}
                {success && <div className="success-message">{success}</div>}
                <div className="btn">
                    <button type="submit" className="submit-button">Сохранить изменения</button>
                </div>
            </form>
        </div>
    );
};

export default ProfileEditDate;