import { useState, useEffect } from "react";
import "./check_course_component.css";
import { PDFDocument, PDFTextField } from "pdf-lib";
import fontkit from "@pdf-lib/fontkit";
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
            const dataResponse = await fetch(`http://localhost:5000/user/contract-data/${courseId}`, {
                method: "GET",
                credentials: "include",
            });
            if (!dataResponse.ok) {
                throw new Error("Ошибка при получении данных договора");
            }
            const formData = await dataResponse.json();

            const pdfResponse = await fetch("http://localhost:5000/contract-template");
            if (!pdfResponse.ok) {
                throw new Error("Ошибка при загрузке шаблона PDF");
            }
            const pdfBytes = await pdfResponse.arrayBuffer();

            const pdfDoc = await PDFDocument.load(pdfBytes);
            const form = pdfDoc.getForm();

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

            const fieldNames = form.getFields().map(field => field.getName());
            console.log("Доступные поля формы:", fieldNames);

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
                DocumentDay: ["DocumentDay"],
                CourseEnd: ["CourseEnd"],
                Time: ["Time"],
                TimeDayStart: ["TimeDayStart"],
                TimeMonthStart: ["TimeMonthStart"],
                TimeDayEnd: ["TimeDayEnd"],
                TimeMonthEnd: ["TimeMonthEnd"],
                HowGive: ["HowGive"],
            };

            for (const [key, value] of Object.entries(formData)) {
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

            form.getFields().forEach(field => {
                if (field instanceof PDFTextField) {
                    try {
                        field.updateAppearances(font);
                    } catch (appearanceError) {
                        console.warn(`Ошибка при обновлении внешнего вида поля ${field.getName()}: ${appearanceError.message}`);
                    }
                }
            });

            form.flatten();
            const updatedPdfBytes = await pdfDoc.save();

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
        } catch (err) {
            console.error("Ошибка:", err);
            alert(err.message);
        }
    };

    const handleDownloadDocument = async (courseId, status) => {
        try {
            const dataResponse = await fetch(`http://localhost:5000/user/download-document/${courseId}`, {
                method: "GET",
                credentials: "include",
            });
            if (!dataResponse.ok) {
                throw new Error("Ошибка при получении данных документа");
            }
            const formData = await dataResponse.json();

            const pdfResponse = await fetch(`http://localhost:5000/user/document-template/${courseId}`);
            if (!pdfResponse.ok) {
                throw new Error("Ошибка при загрузке шаблона PDF");
            }
            const pdfBytes = await pdfResponse.arrayBuffer();

            const pdfDoc = await PDFDocument.load(pdfBytes);
            const form = pdfDoc.getForm();

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

            const fieldNames = form.getFields().map(field => field.getName());
            console.log("Доступные поля формы:", fieldNames);

            const fieldMap = {
                Full_Name: ["Full_Name"],
                Course_Title: ["Course_Title"],
                Start_course: ["Start_course"],
                End_course: ["End_course"],
                Date: ["Date"],
                Random_number: ["Random_number"],
            };

            for (const [key, value] of Object.entries(formData)) {
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

            form.getFields().forEach(field => {
                if (field instanceof PDFTextField) {
                    try {
                        field.updateAppearances(font);
                    } catch (appearanceError) {
                        console.warn(`Ошибка при обновлении внешнего вида поля ${field.getName()}: ${appearanceError.message}`);
                    }
                }
            });

            form.flatten();
            const updatedPdfBytes = await pdfDoc.save();

            const docType = status === "Завершил" ? "диплом" : "сертификат";
            const blob = new Blob([updatedPdfBytes], { type: "application/pdf" });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            a.download = `${docType}_${courseId}.pdf`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);

            alert(`${docType.charAt(0).toUpperCase() + docType.slice(1)} успешно скачан!`);
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
                <button
                    onClick={() => setSelectedStatus("Завершил")}
                    className={`prifle_nav_button ${selectedStatus === "Завершил" ? "active" : ""}`}
                >
                    Завершённые
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
                                <div className="course-items_1">
                                    <span className="course-title">Название: {course.courseTitle}</span>
                                    <span className="course-status">Статус: {course.status}</span>
                                    {course.groupId && (
                                        <span className="course-group">
                                            Группа: {course.groupName}
                                        </span>
                                    )}
                                </div>
                                <div className="course-items">

                                    {course.rejectReason && (
                                        <span className="reject-reason">
                                            Причина отказа: {course.rejectReason}
                                        </span>
                                    )}
                                    {course.status === "Ожидание" && (
                                        <button
                                            className="reject-btn"
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
                                            onClick={() => handleOpenPaymentModal(course)}
                                            className="download-contract-button"
                                        >
                                            Оплатить
                                        </button>
                                    )}
                                    {(course.status === "Завершил" || course.status === "Отчисленный") && (
                                        <button
                                            className="download-document-button"
                                            onClick={() => handleDownloadDocument(course._id, course.status)}
                                        >
                                            {course.status === "Завершил" ? "Скачать диплом" : "Скачать сертификат"}
                                        </button>
                                    )}
                                </div>

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