import React, { useState, useEffect } from 'react';
import './profile_edit_date.css';

const ProfileEditDate = () => {
    const [userData, setUserData] = useState({
        lastname: '',
        firstname: '',
        middlename: '',
        education: '',
        phone: '',
        birthdate: '',
        birthplace: '',
        homeaddress: '',
        workplace: '',
        jobtitle: '',
        oldPassword: '',
        newPassword: '',
        confirmPassword: '',
        agreetoprocessing: false,
    });

    const [educations, setEducations] = useState([]);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [addressSuggestions, setAddressSuggestions] = useState([]);
    const [showSuggestions, setShowSuggestions] = useState(false);

    // Загрузка данных пользователя и списка образований
    useEffect(() => {
        const fetchData = async () => {
            try {
                // Загрузка данных пользователя
                const userResponse = await fetch('http://localhost:5000/profile', {
                    credentials: 'include',
                });
                if (userResponse.ok) {
                    const userData = await userResponse.json();
                    setUserData(prev => ({
                        ...prev,
                        ...userData,
                        educationId: userData.educationId?._id || userData.educationId || ''
                    }));
                }

                // Загрузка списка образований
                const eduResponse = await fetch('http://localhost:5000/admin/educations', {
                    credentials: 'include',
                });
                if (eduResponse.ok) {
                    const eduData = await eduResponse.json();
                    setEducations(eduData || []);
                }
            } catch (err) {
                console.error('Ошибка загрузки:', err);
                setError('Ошибка при загрузке данных');
            }
        };

        fetchData();
    }, []);

    // Функция для запроса подсказок адреса через DaData API
    const fetchAddressSuggestions = async (query) => {
        try {
            const response = await fetch("http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Accept": "application/json",
                    "Authorization": "7a08000a1dd10d29451df5dfd0295c8076c7ad89"
                },
                body: JSON.stringify({
                    query: query,
                    count: 5,
                    locations: [{ country: "*" }] // Можно уточнить поиск по стране/региону
                }),
            });

            const data = await response.json();
            setAddressSuggestions(data.suggestions || []);
        } catch (err) {
            console.error("Ошибка при получении подсказок адреса:", err);
            setAddressSuggestions([]);
        }
    };

    // Обработчик изменения адреса
    const handleAddressChange = (e) => {
        const { value } = e.target;
        setUserData(prev => ({ ...prev, homeaddress: value }));

        if (value.length > 2) {
            fetchAddressSuggestions(value);
            setShowSuggestions(true);
        } else {
            setShowSuggestions(false);
        }
    };

    // Выбор подсказки адреса
    const selectAddressSuggestion = (suggestion) => {
        setUserData(prev => ({
            ...prev,
            homeaddress: suggestion.value
        }));
        setShowSuggestions(false);
    };

    // Валидация кириллицы
    const validateCyrillic = (text) => {
        const cyrillicRegex = /^[а-яА-ЯёЁ\s-]+$/;
        return cyrillicRegex.test(text);
    };

    // Форматирование первой буквы в заглавную
    const capitalizeFirstLetter = (text) => {
        return text
            .split(' ')
            .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
            .join(' ');
    };

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;

        // Валидация для ФИО
        if (name === 'lastname' || name === 'firstname' || name === 'middlename') {
            if (!validateCyrillic(value) && value !== '') {
                setError('ФИО должно содержать только кириллицу');
                return;
            }
        }

        // Форматирование ФИО
        let formattedValue = value;
        if (name === 'lastname' || name === 'firstname' || name === 'middlename') {
            formattedValue = capitalizeFirstLetter(value);
        }

        setUserData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : formattedValue,
        }));
        setError('');
    };

    // Валидация и форматирование телефона
    const validatePhone = (phone) => {
        const phoneRegex = /^\+7 \(\d{3}\) \d{3}-\d{2}-\d{2}$/;
        return phoneRegex.test(phone);
    };

    const handlePhoneChange = (e) => {
        const { value } = e.target;
        let formattedValue = value.replace(/\D/g, '');

        if (formattedValue.length > 0) {
            formattedValue = `+7 (${formattedValue.substring(1, 4)}) ${formattedValue.substring(4, 7)}-${formattedValue.substring(7, 9)}-${formattedValue.substring(9, 11)}`;
        }

        setUserData(prev => ({
            ...prev,
            phone: formattedValue,
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        // Валидация обязательных полей
        if (!userData.lastname || !userData.firstname || !userData.birthdate || !userData.agreetoprocessing) {
            setError('Пожалуйста, заполните все обязательные поля');
            return;
        }

        // Валидация телефона
        if (!validatePhone(userData.phone)) {
            setError('Номер телефона должен быть в формате +7 (XXX) XXX-XX-XX');
            return;
        }

        if (!userData.agreetoprocessing) {
            setError('Необходимо согласие на обработку данных');
            return;
        }

        if (userData.newPassword && userData.newPassword !== userData.confirmPassword) {
            setError('Новый пароль и подтверждение пароля не совпадают');
            return;
        }

        // Подготовка данных для отправки
        const updateData = {
            lastname: userData.lastname,
            firstname: userData.firstname,
            middlename: userData.middlename,
            education: userData.education,
            phone: userData.phone,
            birthdate: userData.birthdate,
            birthplace: userData.birthplace,
            homeaddress: userData.homeaddress,
            workplace: userData.workplace,
            jobtitle: userData.jobtitle,
            agreetoprocessing: userData.agreetoprocessing,
        };

        // Добавляем пароль, если он изменен
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

    // Получаем текущую дату для ограничения выбора даты рождения
    const today = new Date().toISOString().split('T')[0];

    return (
        <div className="profile-edit-container">
            <form onSubmit={handleSubmit}>
                <div className="from_input_group">
                    <div className="input_group">
                        <div className="form-group">
                            <label>Фамилия:</label>
                            <input
                                type="text"
                                name="lastname"
                                value={userData.lastname}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Имя:</label>
                            <input
                                type="text"
                                name="firstname"
                                value={userData.firstname}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Отчество:</label>
                            <input
                                type="text"
                                name="middlename"
                                value={userData.middlename}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>Дата рождения:</label>
                            <input
                                type="date"
                                name="birthdate"
                                value={userData.birthdate}
                                onChange={handleChange}
                                max={today}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Место рождения:</label>
                            <input
                                type="text"
                                name="birthplace"
                                value={userData.birthplace}
                                onChange={handleChange}
                            />
                        </div>
                    </div>
                    <div className="input_group">
                        <div className="form-group">
                            <label>Образование:</label>
                            <select
                                name="educationId"
                                value={userData.educationId || ''}
                                onChange={(e) => setUserData({
                                    ...userData,
                                    educationId: e.target.value
                                })}
                            >
                                <option value="">Выберите образование</option>
                                {educations.map(edu => (
                                    <option key={edu._id} value={edu._id}>
                                        {edu.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div className="form-group">
                            <label>Место работы и должность:</label>
                            <div className="form-group_input">
                                <input
                                    placeholder='Место работы'
                                    type="text"
                                    name="workplace"
                                    value={userData.workplace}
                                    onChange={handleChange}
                                />
                                <input
                                    placeholder='Должность'
                                    type="text"
                                    name="jobtitle"
                                    value={userData.jobtitle}
                                    onChange={handleChange}
                                />
                            </div>
                        </div>
                        <div className="form-group">
                            <label>Номер телефона:</label>
                            <input
                                type="text"
                                name="phone"
                                value={userData.phone}
                                onChange={handlePhoneChange}
                                placeholder="+7 (XXX) XXX-XX-XX"
                            />
                        </div>
                        <div className="form-group">
                            <label>Домашний адрес (прописка):</label>
                            <input
                                type="text"
                                name="homeaddress"
                                value={userData.homeaddress}
                                onChange={handleAddressChange}
                                placeholder="Начните вводить адрес"
                            />
                            {showSuggestions && addressSuggestions.length > 0 && (
                                <div className="address-suggestions">
                                    {addressSuggestions.map((suggestion, index) => (
                                        <div
                                            key={index}
                                            className="suggestion-item"
                                            onClick={() => selectAddressSuggestion(suggestion)}
                                        >
                                            {suggestion.value}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    name="agreetoprocessing"
                                    checked={userData.agreetoprocessing}
                                    onChange={handleChange}
                                    required
                                />
                                Согласен на обработку персональных данных
                            </label>
                        </div>
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
                        <button
                            type="button"
                            onClick={toggleShowPassword}
                            className="show-password-btn"
                        >
                            {showPassword ? "Скрыть пароль" : "Показать пароль"}
                        </button> */}
                    </div>
                </div>
                {error && <div className="error-message">{error}</div>}
                {success && <div className="success-message">{success}</div>}
                <div className="btn">
                    <button type="submit" className="submit-button">
                        Сохранить изменения
                    </button>
                </div>
            </form>
        </div>
    );
};

export default ProfileEditDate;