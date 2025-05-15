import { useState, useEffect } from "react";
import "./check_course_component.css";
import { PDFDocument, StandardFonts, rgb } from "pdf-lib";
import qr from "../../../assets/mustafin_qr.jpg";
import ArialFont from "../../../assets/fonts/Arial.ttf";

const CheckCourse = () => {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedStatus, setSelectedStatus] = useState("all");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedCourse, setSelectedCourse] = useState(null);

    const fetchCourses = async (status) => {
        try {
            let endpoint = "/user/courses";
            if (status !== "all") {
                endpoint = `/user/courses/status?status=${status}`;
            }
            const response = await fetch(`http://localhost:5000${endpoint}`, {
                method: "GET",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при загрузке курсов");
            }
            const data = await response.json();
            if (!Array.isArray(data)) {
                console.error("Сервер вернул некорректные данные:", data);
                setCourses([]);
                return;
            }
            setCourses(data);
        } catch (err) {
            setError(err.message);
            setCourses([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCourses(selectedStatus);
    }, [selectedStatus]);

    const handleDownloadContract = async (courseId) => {
        try {
            // 1. Получить данные договора
            const dataResponse = await fetch(`http://localhost:5000/user/contract-data/${courseId}`, {
                method: "GET",
                credentials: "include",
            });
            if (!dataResponse.ok) {
                throw new Error("Ошибка при получении данных договора");
            }
            const formData = await dataResponse.json();

            // 2. Загрузить шаблон PDF
            const pdfResponse = await fetch("http://localhost:5000/contract-template");
            if (!pdfResponse.ok) {
                throw new Error("Ошибка при загрузке шаблона PDF");
            }
            const pdfBytes = await pdfResponse.arrayBuffer();

            // 3. Загрузить PDF в pdf-lib
            const pdfDoc = await PDFDocument.load(pdfBytes);
            const form = pdfDoc.getForm();

            // 4. Загрузить шрифт Arial
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
                console.warn("Не удалось загрузить Arial, использую Helvetica:", fontError);
                font = await pdfDoc.embedFont(StandardFonts.Helvetica); // Helvetica как запасной вариант
            }
            // 5. Заполнить поля формы
            try {
                // Логируем доступные поля для отладки
                const fieldNames = form.getFields().map(field => field.getName());
                console.log("Доступные поля формы:", fieldNames);

                // Карта полей (основана на полях из PDF)
                const fieldMap = {
                    FullName: ["FullName"],
                    CourseTitle: ["CourseTitle"],
                    CourseDuration: ["CourseDuration"],
                    CoursePrice: ["CoursePrice"],
                    DocumentId: ["DocumentId"],
                    Adress: ["Adress"],
                    PasportDate: ["PasportDate"],
                    SNILS: ["SNILS"],
                    Phone: ["Phone"],
                    Email: ["Email"],
                    DocumentDay: ["DocumentDay"], // Новое поле
                    CourseEnd: ["CourseEnd"],     // Новое поле
                    Time: ["Time"],               // Новое поле
                    TimeDayStart: ["TimeDayStart"],
                    TimeMonthStart: ["TimeMonthStart"],
                    TimeDayEnd: ["TimeDayEnd"],
                    TimeMonthEnd: ["TimeMonthEnd"],
                    HowGive: ["HowGive"],         // Новое поле
                };
                for (const [key, value] of Object.entries(formData)) {
                    const possibleNames = fieldMap[key] || [key];
                    let fieldFilled = false;
                    for (const fieldName of possibleNames) {
                        try {
                            const field = form.getTextField(fieldName);
                            if (field) {
                                field.setText(value);
                                // Обновляем внешний вид поля с кастомным шрифтом
                                field.updateAppearances(font);
                                field.enableMultiline(); // На случай длинного текста
                                console.log(`Заполнено поле ${fieldName}: ${value}`);
                                fieldFilled = true;
                                break;
                            }
                        } catch (fieldError) {
                            console.warn(`Ошибка при заполнении поля ${fieldName}:`, fieldError.message);
                        }
                    }
                    if (!fieldFilled) {
                        console.warn(`Поле для ключа ${key} не найдено среди ${possibleNames.join(", ")}`);
                    }
                }

                // 6. Сохранить PDF
                form.flatten(); // Зафиксировать поля
                const updatedPdfBytes = await pdfDoc.save();

                // 7. Скачать PDF
                const blob = new Blob([updatedPdfBytes], { type: "application/pdf" });
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.style.display = "none";
                a.href = url;
                a.download = `contract_${courseId}.pdf`;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);

                alert("Договор успешно скачан! Заполните его и отправьте нам!");
            } catch (formError) {
                throw new Error(`Ошибка при обработке формы PDF: ${formError.message}`);
            }
        } catch (err) {
            console.error("Ошибка:", err);
            alert(err.message);
        }
    };

    const handleUploadContract = async (courseId, file) => {
        const formData = new FormData();
        formData.append("contract", file);
        try {
            const response = await fetch(`http://localhost:5000/user/upload-contract/${courseId}`, {
                method: "POST",
                credentials: "include",
                body: formData,
            });
            if (!response.ok) {
                throw new Error("Ошибка при загрузке договора");
            }
            const data = await response.json();
            alert(data.message || "Договор успешно загружен!");
            setCourses((prevCourses) =>
                prevCourses.map((course) =>
                    course._id === courseId ? { ...course, contractUploaded: true } : course
                )
            );
        } catch (err) {
            alert(err.message);
        }
    };

    const handleWithdrawRegistration = async (registrationId) => {
        try {
            const response = await fetch(`http://localhost:5000/user/withdraw-registration/${registrationId}`, {
                method: "DELETE",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при отзыве заявки");
            }
            const data = await response.json();
            if (data.success) {
                setCourses((prevCourses) => prevCourses.filter(course => course._id !== registrationId));
                alert("Заявка успешно отозвана!");
            } else {
                throw new Error("Ошибка при отзыве заявки");
            }
        } catch (err) {
            alert(err.message);
        }
    };

    const handleOpenPaymentModal = (course) => {
        setSelectedCourse(course);
        setIsModalOpen(true);
    };

    const handleClosePaymentModal = () => {
        setIsModalOpen(false);
        setSelectedCourse(null);
    };

    const handlePayment = async () => {
        try {
            const response = await fetch(`http://localhost:5000/user/pay-course/${selectedCourse._id}`, {
                method: "POST",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при оплате курса");
            }
            const data = await response.json();
            if (data.success) {
                setCourses((prevCourses) =>
                    prevCourses.map((course) =>
                        course._id === selectedCourse._id ? { ...course, status: "Оплаченный" } : course
                    )
                );
                alert("Курс успешно оплачен!");
                handleClosePaymentModal();
            } else {
                throw new Error("Ошибка при оплате курса");
            }
        } catch (err) {
            alert(err.message);
        }
    };

    return (
        <div className="check-course-container">
            <div className="profile_course_navs">
                <button
                    onClick={() => setSelectedStatus("all")}
                    className={`prifle_nav_button ${selectedStatus === "all" ? "active" : ""}`}
                >
                    Все курсы
                </button>
                <button
                    onClick={() => setSelectedStatus("Ожидание")}
                    className={`prifle_nav_button ${selectedStatus === "Ожидание" ? "active" : ""}`}
                >
                    Ожидают одобрения
                </button>
                <button
                    onClick={() => setSelectedStatus("Одобренный")}
                    className={`prifle_nav_button ${selectedStatus === "Одобренный" ? "active" : ""}`}
                >
                    Одобренные
                </button>
                <button
                    onClick={() => setSelectedStatus("Принят")}
                    className={`prifle_nav_button ${selectedStatus === "Принят" ? "active" : ""}`}
                >
                    Принятые
                </button>
                <button
                    onClick={() => setSelectedStatus("Отклоненный")}
                    className={`prifle_nav_button ${selectedStatus === "Отклоненный" ? "active" : ""}`}
                >
                    Отклоненные
                </button>
                <button
                    onClick={() => setSelectedStatus("Отчисленный")}
                    className={`prifle_nav_button ${selectedStatus === "Отчисленный" ? "active" : ""}`}
                >
                    Отчисленный
                </button>
            </div>

            <div className="course-list-container">
                <h2 className="course-list-title">Список курсов:</h2>
                {loading && <p className="loading">Загрузка...</p>}
                {error && <p className="error">{error}</p>}
                {!loading && !error && courses.length === 0 && (
                    <p className="no-courses-message">Нет доступных курсов для выбранного статуса.</p>
                )}
                {!loading && !error && courses.length > 0 && (
                    <ul className="course-list">
                        {courses.map((course, index) => (
                            <li key={index} className="course-item">
                                <span className="course-title">Название: {course.courseTitle}</span>
                                <span className="course-status">Статус: {course.status}</span>
                                {course.groupId && (
                                    <span className="course-group">
                                        Группа: {course.groupName}
                                    </span>
                                )}
                                {course.rejectReason && (
                                    <span className="reject-reason">
                                        Причина отказа: {course.rejectReason}
                                    </span>
                                )}
                                {course.status === "Ожидание" && (
                                    <button
                                        className="withdraw-button"
                                        onClick={() => handleWithdrawRegistration(course._id)}
                                    >
                                        Отозвать заявку
                                    </button>
                                )}
                                {course.status === "Одобренный" && (
                                    <>
                                        {!course.contractUploaded && (
                                            <div className="download-buttons">
                                                <button
                                                    className="download-contract-button"
                                                    onClick={() => handleDownloadContract(course._id)}
                                                >
                                                    Скачать договор
                                                </button>
                                                <label htmlFor={`upload-${course._id}`} className="upload-label">
                                                    Загрузить файл
                                                </label>
                                                <input
                                                    id={`upload-${course._id}`}
                                                    type="file"
                                                    accept=".pdf"
                                                    style={{ display: "none" }}
                                                    onChange={(e) => {
                                                        const file = e.target.files[0];
                                                        if (file) {
                                                            handleUploadContract(course._id, file);
                                                        }
                                                    }}
                                                />
                                            </div>
                                        )}
                                        <span className="file-status">
                                            {course.contractUploaded
                                                ? "Файл успешно загружен."
                                                : "Файла на проверку нет."}
                                        </span>
                                    </>
                                )}
                                {course.status === "Принят" && (
                                    <button
                                        className="pay-button"
                                        onClick={() => handleOpenPaymentModal(course)}
                                    >
                                        Оплатить
                                    </button>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            {isModalOpen && (
                <div className="payment-modal">
                    <div className="payment-modal-content">
                        <h2>Оплата курса: {selectedCourse.courseTitle}</h2>
                        <div className="qr-code-placeholder">
                            <img src={qr} alt="QR Code" />
                        </div>
                        <div className="payment-modal-buttons">
                            <button onClick={handlePayment}>Оплатил</button>
                            <button onClick={handleClosePaymentModal}>Назад</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default CheckCourse;