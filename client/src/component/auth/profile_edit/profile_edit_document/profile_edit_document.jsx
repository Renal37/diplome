import React, { useState, useEffect } from 'react';
import './profile_edit_document.css';

const ProfileEditDocument = () => {
    const [userData, setUserData] = useState({
        passportSeries: '',
        passportNumber: '',
        snils: '',
        agreeToProcessing: false,
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

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
                        passportSeries: data.passportData ? data.passportData.split(' ')[0] : '',
                        passportNumber: data.passportData ? data.passportData.split(' ')[1] : '',
                        snils: data.snils || '',
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

        // Объединяем серию и номер паспорта
        const passportData = `${userData.passportSeries} ${userData.passportNumber}`.trim();

        // Создаем объект для отправки только измененных данных
        const updateData = {
            passportData,
            snils: userData.snils,
            agreeToProcessing: userData.agreeToProcessing,
        };

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

    return (
        <div className="profile-edit-container">
            <form onSubmit={handleSubmit}>
                <div className="from_input_group">
                    <div className="input_group">
                        <div className="form-group">
                            <label>Паспортные данные:</label>
                            <div className="pacport">
                                <input
                                    type="text"
                                    name="passportSeries"
                                    value={userData.passportSeries}
                                    onChange={handleChange}
                                    placeholder="Серия паспорта"
                                />
                                <input
                                    type="text"
                                    name="passportNumber"
                                    value={userData.passportNumber}
                                    onChange={handleChange}
                                    placeholder="Номер паспорта"
                                />
                            </div>
                        </div>
                    </div>
                    <div className="input_group">
                        <div className="form-group">
                            <label>Снилс:</label>
                            <input
                                type="text"
                                name="snils"
                                value={userData.snils}
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

export default ProfileEditDocument;