import React, { useState, useEffect } from 'react';
import './profile_edit_date.css';

const ProfileEditDate = () => {
    const [userData, setUserData] = useState({
        lastName: '',
        firstName: '',
        middleName: '',
        education: '',
        phone: '',
        birthDate: '',
        birthPlace: '',
        homeAddress: '',
        workPlace: '',
        jobTitle: '',
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

                    // Обновляем состояние, сохраняя предыдущие значения
                    setUserData((prevState) => ({
                        ...prevState,
                        lastName: data.lastName || '',
                        firstName: data.firstName || '',
                        middleName: data.middleName || '',
                        education: data.education || '',
                        phone: data.phone || '',
                        birthDate: data.birthDate || '',
                        birthPlace: data.birthPlace || '',
                        homeAddress: data.homeAddress || '',
                        workPlace: data.workPlace || '',
                        jobTitle: data.jobTitle || '',
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

    // Функция для валидации кириллицы
    const validateCyrillic = (text) => {
        const cyrillicRegex = /^[а-яА-ЯёЁ\s-]+$/; // Разрешаем кириллицу, пробелы и дефисы
        return cyrillicRegex.test(text);
    };

    // Функция для форматирования первой буквы в заглавную
    const capitalizeFirstLetter = (text) => {
        return text
            .split(' ') // Разделяем строку на слова
            .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()) // Каждое слово начинаем с заглавной буквы
            .join(' '); // Соединяем слова обратно в строку
    };

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;

        // Валидация для ФИО (только кириллица)
        if (name === 'lastName' || name === 'firstName' || name === 'middleName') {
            if (!validateCyrillic(value) && value !== '') {
                setError('ФИО должно содержать только кириллицу');
                return;
            }
        }

        // Форматирование первой буквы в заглавную для ФИО
        let formattedValue = value;
        if (name === 'lastName' || name === 'firstName' || name === 'middleName') {
            formattedValue = capitalizeFirstLetter(value);
        }

        setUserData((prevState) => ({
            ...prevState,
            [name]: type === 'checkbox' ? checked : formattedValue,
        }));
        setError(''); // Сбрасываем ошибку при успешном вводе
    };

    const validatePhone = (phone) => {
        const phoneRegex = /^\+7 \(\d{3}\) \d{3}-\d{2}-\d{2}$/;
        return phoneRegex.test(phone);
    };

    const handlePhoneChange = (e) => {
        const { value } = e.target;
        let formattedValue = value.replace(/\D/g, ''); // Удаляем все нецифровые символы

        if (formattedValue.length > 0) {
            formattedValue = `+7 (${formattedValue.substring(1, 4)}) ${formattedValue.substring(4, 7)}-${formattedValue.substring(7, 9)}-${formattedValue.substring(9, 11)}`;
        }

        setUserData((prevState) => ({
            ...prevState,
            phone: formattedValue,
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccess('');

        // Валидация обязательных полей
        if (!userData.lastName || !userData.firstName || !userData.birthDate || !userData.agreeToProcessing) {
            setError('Пожалуйста, заполните все обязательные поля');
            return;
        }

        // Валидация номера телефона
        if (!validatePhone(userData.phone)) {
            setError('Номер телефона должен быть в формате +7 (XXX) XXX-XX-XX');
            return;
        }

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
            phone: userData.phone,
            birthDate: userData.birthDate,
            birthPlace: userData.birthPlace,
            homeAddress: userData.homeAddress,
            workPlace: userData.workPlace,
            jobTitle: userData.jobTitle,
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
                                name="lastName"
                                value={userData.lastName}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Имя:</label>
                            <input
                                type="text"
                                name="firstName"
                                value={userData.firstName}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Отчество:</label>
                            <input
                                type="text"
                                name="middleName"
                                value={userData.middleName}
                                onChange={handleChange}
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Дата рождения:</label>
                            <input
                                type="date"
                                name="birthDate"
                                value={userData.birthDate}
                                onChange={handleChange}
                                max={today} // Ограничение на выбор даты
                                required
                            />
                        </div>
                        <div className="form-group">
                            <label>Место рождения:</label>
                            <input
                                type="text"
                                name="birthPlace"
                                value={userData.birthPlace}
                                onChange={handleChange}
                            />
                        </div>
                    </div>
                    <div className="input_group">
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
                        <div className="form-group">
                            <label>Место работы и занимая должность:</label>
                            <div className="form-group_input">
                                <input
                                    placeholder='Место работы'
                                    type="text"
                                    name="workPlace"
                                    value={userData.workPlace}
                                    onChange={handleChange}
                                />
                                <input
                                    placeholder='Должность'
                                    type="text"
                                    name="jobTitle"
                                    value={userData.jobTitle}
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
                                name="homeAddress"
                                value={userData.homeAddress}
                                onChange={handleChange}
                            />
                        </div>
                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    name="agreeToProcessing"
                                    checked={userData.agreeToProcessing}
                                    onChange={handleChange}
                                    required
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