import React, { useState, useEffect } from "react";
import "./profile_edit_document.css";
import { PDFDocument, rgb, PDFTextField } from "pdf-lib";
import fontkit from "@pdf-lib/fontkit";
import ArialFont from "../../../../assets/fonts/Arial.ttf";

const formatSnils = (input) => {
    const numbers = input.replace(/\D/g, "");
    const part1 = numbers.slice(0, 3);
    const part2 = numbers.slice(3, 6);
    const part3 = numbers.slice(6, 9);
    const part4 = numbers.slice(9, 11);
    return `${part1}-${part2}-${part3} ${part4}`.trim();
};

const formatDate = (dateStr) => {
    if (!dateStr) return { day: "", month: "", year: "" };
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return { day: "", month: "", year: "" }; // Проверка на валидность даты
    const day = date.getDate().toString().padStart(2, "0");
    const year = date.getFullYear().toString(); // Полный год для GiveYear
    const monthNames = [
        "января", "февраля", "марта", "апреля", "мая", "июня",
        "июля", "августа", "сентября", "октября", "ноября", "декабря"
    ];
    const month = monthNames[date.getMonth()] || "";
    return { day, month, year };
};

const getCurrentDate = () => {
    const date = new Date();
    const day = date.getDate().toString().padStart(2, "0");
    const year = date.getFullYear().toString().slice(-2); // Последние две цифры для NewYear
    const monthNames = [
        "января", "февраля", "марта", "апреля", "мая", "июня",
        "июля", "августа", "сентября", "октября", "ноября", "декабря"
    ];
    const month = monthNames[date.getMonth()];
    return { day, month, year };
};

const ProfileEditDocument = () => {
    const [userData, setUserData] = useState({
        passportSeries: "",
        passportNumber: "",
        passportIssuedBy: "",
        passportIssueDate: "",
        snils: "",
        agreetoprocessing: false,
    });
    const [profile, setProfile] = useState(null);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

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
                setUserData({
                    passportSeries: data.passportdata ? data.passportdata.split(' ')[0] : '',
                    passportNumber: data.passportdata ? data.passportdata.split(' ')[1] : '',
                    passportIssuedBy: data.passportissuedby || '',
                    passportIssueDate: data.passportissuedate || '',
                    snils: data.snils || '',
                    agreetoprocessing: data.agreetoprocessing || false,
                });
            } catch (err) {
                setError(err.message);
            }
        };
        fetchProfile();
    }, []);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        if (name === "snils") {
            const numbers = value.replace(/\D/g, "");
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

        if (!userData.snils.match(/^\d{3}-\d{3}-\d{3} \d{2}$/)) {
            setError("Неверный формат СНИЛСа (XXX-XXX-XXX XX)");
            return;
        }

        if (!userData.agreetoprocessing) {
            setError("Необходимо согласие на обработку данных");
            return;
        }

        if (userData.passportIssueDate) {
            const selectedDate = new Date(userData.passportIssueDate);
            const today = new Date();
            if (selectedDate > today) {
                setError("Дата выдачи паспорта не может быть в будущем");
                return;
            }
            if (selectedDate < new Date("1900-01-01")) {
                setError("Дата выдачи паспорта слишком старая");
                return;
            }
        }

        const updateData = {
            passportData: `${userData.passportSeries} ${userData.passportNumber}`.trim(),
            passportIssuedBy: userData.passportIssuedBy,
            passportIssueDate: userData.passportIssueDate,
            snils: userData.snils,
            agreeToProcessing: userData.agreetoprocessing,
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

    const handleDownloadContract = async () => {
        try {
            // 1. Формируем данные для заполнения из profile и userData
            if (!profile) {
                throw new Error("Профиль пользователя не загружен");
            }

            const { day: giveDay, month: giveMounth, year: giveYear } = formatDate(userData.passportIssueDate);
            const { day: newDay, month: newMounth, year: newYear } = getCurrentDate();

            const formData = {
                FullName: `${profile.lastname || ''} ${profile.firstname || ''} ${profile.middlename || ''}`.trim(),
                Adress: profile.homeaddress || '',
                PassportSeria: userData.passportSeries || '',
                PassportNumber: userData.passportNumber || '',
                GiveDay: giveDay,
                GiveMounth: giveMounth,
                GiveYear: giveYear,
                HowGive: userData.passportIssuedBy || '',
                Phone: profile.phone || '',
                NewDay: newDay,
                NewMounth: newMounth,
                NewYear: newYear,
            };

            // Логируем formData для отладки
            console.log("formData:", formData);

            // 2. Загрузить шаблон PDF
            const pdfResponse = await fetch("http://localhost:5000/user/download-document", {
                method: "GET",
                credentials: "include",
            });
            if (!pdfResponse.ok) {
                throw new Error("Ошибка при загрузке шаблона PDF");
            }
            const pdfBytes = await pdfResponse.arrayBuffer();

            // 3. Загрузить PDF в pdf-lib
            const pdfDoc = await PDFDocument.load(pdfBytes);
            const form = pdfDoc.getForm();

            // 4. Регистрация fontkit и загрузка шрифта Arial
            pdfDoc.registerFontkit(fontkit);
            let font;
            try {
                const fontResponse = await fetch(ArialFont);
                if (!fontResponse.ok) {
                    throw new Error("Не удалось загрузить шрифт Arial");
                }
                const fontBytes = await fontResponse.arrayBuffer();
                font = await pdfDoc.embedFont(fontBytes, { subset: true });
                console.log("Шрифт Arial успешно загружен");
            } catch (fontError) {
                throw new Error(`Не удалось загрузить шрифт Arial: ${fontError.message}`);
            }

            // 5. Заполнить поля формы
            try {
                const fieldNames = form.getFields().map(field => field.getName());
                console.log("Доступные поля формы:", fieldNames);

                const fieldMap = {
                    FullName: ["FullName"],
                    Adress: ["Adress"],
                    PassportSeria: ["PassportSeria"],
                    PassportNumber: ["PassportNumber"],
                    GiveDay: ["GiveDay"],
                    GiveMounth: ["GiveMounth"],
                    GiveYear: ["GiveYear"],
                    HowGive: ["HowGive"],
                    Phone: ["Phone"],
                    NewDay: ["NewDay"],
                    NewMounth: ["NewMounth"],
                    NewYear: ["NewYear"],
                };

                for (const [key, value] of Object.entries(formData)) {
                    if (typeof value !== "string" || value == null || value === "") {
                        console.warn(`Некорректное или пустое значение для ключа ${key}: ${value}`);
                        continue;
                    }
                    const possibleNames = fieldMap[key] || [key];
                    let fieldFilled = false;
                    for (const fieldName of possibleNames) {
                        try {
                            const field = form.getTextField(fieldName);
                            if (field) {
                                console.log(`Попытка заполнить поле ${fieldName} значением: ${value}`);
                                field.setText(value);
                                field.updateAppearances(font);
                                field.enableMultiline();
                                console.log(`Успешно заполнено поле ${fieldName}: ${value}`);
                                fieldFilled = true;
                                break;
                            } else {
                                console.warn(`Поле ${fieldName} не найдено в форме`);
                            }
                        } catch (fieldError) {
                            console.warn(`Ошибка при заполнении поля ${fieldName}: ${fieldError.message}`);
                        }
                    }
                    if (!fieldFilled) {
                        console.warn(`Поле для ключа ${key} не найдено среди ${possibleNames.join(", ")}`);
                    }
                }

                // Применяем шрифт ко всем текстовым полям
                form.getFields().forEach(field => {
                    if (field instanceof PDFTextField) {
                        try {
                            field.updateAppearances(font);
                        } catch (appearanceError) {
                            console.warn(`Ошибка при обновлении внешнего вида поля ${field.getName()}: ${appearanceError.message}`);
                        }
                    }
                });

                // 6. Сохранить PDF
                form.flatten();
                const updatedPdfBytes = await pdfDoc.save();

                // 7. Скачать PDF
                const blob = new Blob([updatedPdfBytes], { type: "application/pdf" });
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.style.display = "none";
                a.href = url;
                a.download = "согласие_на_обработку_данных.pdf";
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);

                alert("Согласие успешно скачано! Проверьте заполненные данные и отправьте нам!");
            } catch (formError) {
                throw new Error(`Ошибка при обработке формы PDF: ${formError.message}`);
            }
        } catch (err) {
            console.error("Ошибка:", err);
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
                throw new Error("Ошибка при загрузке согласия");
            }
            alert("Согласие успешно загружено!");
        } catch (err) {
            alert(err.message);
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
                            <input
                                type="text"
                                name="passportIssuedBy"
                                value={userData.passportIssuedBy}
                                onChange={handleChange}
                                placeholder="Кем выдан паспорт"
                                maxLength={100}
                            />
                            <input
                                type="date"
                                name="passportIssueDate"
                                value={userData.passportIssueDate}
                                onChange={handleChange}
                                placeholder="Дата выдачи паспорта"
                            />
                        </div>
                    </div>
                    <div className="input_group">
                        <div className="form-group">
                            <label>СНИЛС:</label>
                            <input
                                type="text"
                                name="snils"
                                value={userData.snils}
                                onChange={handleChange}
                                placeholder="Формат: XXX-XXX-XXX XX"
                                maxLength={14}
                                pattern="\d{3}-\d{3}-\d{3} \d{2}"
                            />
                        </div>
                        <div className="form-group">
                            <label>
                                <input
                                    type="checkbox"
                                    name="agreetoprocessing"
                                    checked={userData.agreetoprocessing}
                                    onChange={handleChange}
                                />
                                Согласен на обработку персональных данных
                            </label>
                        </div>
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
                        onClick={() => handleDownloadContract()}
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
                                handleUploadContract(profile._id, file);
                            }
                        }}
                    />
                </div>
            )}
        </div>
    );
};

export default ProfileEditDocument;