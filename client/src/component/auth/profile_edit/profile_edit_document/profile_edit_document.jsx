import React, { useState, useEffect } from "react";
import "./profile_edit_document.css";

const formatSnils = (input) => {
    // Оставляем только цифры
    const numbers = input.replace(/\D/g, "");
    // Форматируем в виде XXX-XXX-XXX XX
    const part1 = numbers.slice(0, 3);
    const part2 = numbers.slice(3, 6);
    const part3 = numbers.slice(6, 9);
    const part4 = numbers.slice(9, 11);

    return `${part1}-${part2}-${part3} ${part4}`.trim();
};

const ProfileEditDocument = () => {
    const [userData, setUserData] = useState({
        passportSeries: "",
        passportNumber: "",
        snils: "",
        agreeToProcessing: false,
    });
    const [profile, setProfile] = useState(null);


    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await fetch("http://localhost:5000/profile", {
                    method: "GET",
                    credentials: "include",
                });

                if (!response.ok) {
                    throw new Error("Ошибка при загрузке профиля");
                }

                const data = await response.json();
                setProfile(data);
            } catch (err) {
                setError(err.message);
            }
        };

        fetchProfile();
    }, []);

    const handleDownloadContract = async () => {
        try {
            const response = await fetch(`http://localhost:5000/user/download-document`, {
                method: "GET",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при скачивании договора");
            }
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            a.download = "согласние.pdf";
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            alert("Вы скачали соглашение, теперь заполните его и отправьте нам!");
        } catch (err) {
            alert(err.message);
        }
    };

    const handleUploadContract = async (userId, file) => {
        const formData = new FormData();
        formData.append("contract", file);

        try {
            const response = await fetch(`http://localhost:5000/user/upload-document/${userId}`, {
                method: "POST",
                body: formData,
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при загрузке договора");
            }
            alert("Договор успешно загружен!");
        } catch (err) {
            alert(err.message);
        }
    };


    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

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

        if (name === "snils") {
            // Оставляем только цифры
            const numbers = value.replace(/\D/g, "");
            // Форматируем СНИЛС при изменении
            const formattedValue = formatSnils(numbers);
            setUserData((prevState) => ({
                ...prevState,
                [name]: formattedValue,
            }));
        } else {
            setUserData((prevState) => ({
                ...prevState,
                [name]: type === "checkbox" ? checked : value,
            }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setSuccess("");

        // Простая валидация
        if (!userData.snils.match(/^\d{3}-\d{3}-\d{3} \d{2}$/)) {
            setError("Неверный формат СНИЛСа (XXX-XXX-XXX XX)");
            return;
        }

        if (!userData.agreeToProcessing) {
            setError("Необходимо согласие на обработку данных");
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
                                    placeholder="Серия паспорта (4 цифры)"
                                    maxLength={4}
                                />
                                <input
                                    type="text"
                                    name="passportNumber"
                                    value={userData.passportNumber}
                                    onChange={handleChange}
                                    placeholder="Номер паспорта (6 цифр)"
                                    maxLength={6}
                                />
                            </div>
                        </div>
                    </div>
                    <div className="input_group">
                        <div className="form-group">
                            <label>СНИЛС:</label>
                            <input
                                type="text" // Используем type="tel" для мобильных устройств
                                name="snils"
                                value={userData.snils}
                                onChange={handleChange}
                                placeholder="Формат: XXX-XXX-XXX XX"
                                maxLength={14} // Ограничиваем длину
                                pattern="\d{3}-\d{3}-\d{3} \d{2}" // Паттерн для валидации
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
                    <button type="submit" className="submit-button">
                        Сохранить изменения
                    </button>
                </div>
            </form>
            {profile && (
                <div className="download-cont">
                    <button
                        className="download-contract-button"
                        onClick={() => handleDownloadContract(profile.ID)}
                    >
                        Скачать согласие
                    </button>
                    <label htmlFor="upload-contract" className="upload-label">
                        Загрузить файл
                    </label>
                    <input
                        id="upload-contract"
                        type="file"
                        accept=".pdf"
                        style={{ display: "none" }}
                        onChange={(e) => {
                            const file = e.target.files[0];
                            if (file) {
                                handleUploadContract(profile.ID, file);
                            }
                        }}
                    />
                </div>
            )}
        </div>
    );
};

export default ProfileEditDocument;