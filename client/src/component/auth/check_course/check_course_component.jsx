import { useState, useEffect } from "react";
import "./check_course_component.css";
import { PDFDocument } from "pdf-lib";

const CheckCourse = () => {
    const [courses, setCourses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [selectedStatus, setSelectedStatus] = useState("all");

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
            const response = await fetch(`http://localhost:5000/user/download-contract/${courseId}`, {
                method: "GET",
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error("Ошибка при скачивании договора");
            }
            const pdfBytes = await response.arrayBuffer();
            const pdfDoc = await PDFDocument.load(pdfBytes);
            const updatedPdfBytes = await pdfDoc.save();
            const blob = new Blob([updatedPdfBytes], { type: "application/pdf" });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.style.display = "none";
            a.href = url;
            a.download = "contract.pdf";
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            alert("Вы скачали договор, теперь заполните его и отправьте нам!");
        } catch (err) {
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
            // Обновляем состояние курсов после успешной загрузки
            setCourses((prevCourses) =>
                prevCourses.map((course) =>
                    course._id === courseId ? { ...course, contractUploaded: false } : course
                )
            );
        } catch (err) {
            alert(err.message);
        }
    };

    return (
        <div className="check-course-container">
            {/* Навигация по статусам курсов */}
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
                    Принят курысы
                </button>
                <button
                    onClick={() => setSelectedStatus("Отклоненный")}
                    className={`prifle_nav_button ${selectedStatus === "Отклоненный" ? "active" : ""}`}
                >
                    Отклоненные
                </button>
            </div>

            {/* Список курсов */}
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
                                {/* Отображение группы, если она есть */}
                                {course.groupId && (
                                    <span className="course-group">
                                        Группа: {course.groupName}
                                    </span>
                                )}
                                {course.status === "Отклоненный" && course.rejectReason && (
                                    <span className="reject-reason">
                                        Причина отказа: {course.rejectReason}
                                    </span>
                                )}
                                {course.status === "Одобренный" && (
                                    <>


                                        {!course.contractUploaded && (
                                            <>
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
                                            </>
                                        )}
                                        <span className="file-status">
                                            {course.contractUploaded
                                                ? "Файл успешно загружен."
                                                : "Файла на проверку нет."}
                                        </span>
                                    </>
                                )}
                            </li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
};

export default CheckCourse;