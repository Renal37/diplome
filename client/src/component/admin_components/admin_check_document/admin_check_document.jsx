import React, { useState, useEffect } from "react";
import "./admin_check_document.css";

const AdminCheckDocument = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [rejectReason, setRejectReason] = useState(""); // Причина отчисления
    const [selectedRegistrationId, setSelectedRegistrationId] = useState(null);
    const [selectedUser, setSelectedUser] = useState(null); // Выбранный пользователь
    const [selectedPdfRegistrationId, setSelectedPdfRegistrationId] = useState(null); // ID заявки с PDF
    const [pdfUrl, setPdfUrl] = useState(""); // URL для просмотра PDF
    const [rejectReasonPdf, setRejectReasonPdf] = useState(""); // Причина отклонения для PDF

    useEffect(() => {
        fetch("http://localhost:5000/admin/course-registrations", {
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setError(data.error);
                } else {
                    const approvedRegistrations = data.filter(
                        (reg) => reg.status === "Одобренный" && reg.contractFilePath && reg.contractFilePath.trim() !== null
                    );
                    setRegistrations(approvedRegistrations);
                    console.log(approvedRegistrations);
                }
                setIsLoading(false);
            })
            .catch((error) => {
                console.error("Error fetching registrations:", error);
                setError("Ошибка при загрузке заявок");
                setIsLoading(false);
            });
    }, []);

    const handleExpel = (registrationId) => {
        setSelectedRegistrationId(registrationId);
    };

    const confirmExpel = () => {
        if (!rejectReason) {
            setError("Укажите причину отчисления");
            return;
        }
        fetch(`http://localhost:5000/admin/expel-registration/${selectedRegistrationId}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ reason: rejectReason }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedRegistrationId
                                ? { ...reg, status: "Отчисленный", rejectReason }
                                : reg
                        )
                    );
                    clearModalState();
                } else {
                    setError(data.message || "Ошибка при отчислении");
                }
            })
            .catch((error) => {
                console.error("Error expelling registration:", error);
                setError("Ошибка при отчислении");
            });
    };





    const handleViewUser = (user) => {
        setSelectedUser(user);
    };

    const handleCloseUserInfo = () => {
        clearModalState();
    };

    const clearModalState = () => {
        setRejectReason("");
        setDocumentType("");
        setSelectedRegistrationId(null);
        setSelectedDocumentId(null);
        setSelectedUser(null);
        setSelectedPdfRegistrationId(null);
        setRejectReasonPdf("");
    };

    const closeModalOnEscape = (e) => {
        if (e.key === "Escape") {
            clearModalState();
        }
    };

    useEffect(() => {
        document.addEventListener("keydown", closeModalOnEscape);
        return () => document.removeEventListener("keydown", closeModalOnEscape);
    }, []);

    const handleApprovePdf = (registrationId) => {
        fetch(`http://localhost:5000/admin/approve-contract/${registrationId}`, {
            method: "POST",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === registrationId
                                ? { ...reg, status: "Принят" }
                                : reg
                        )
                    );
                    setSelectedPdfRegistrationId(null); // Закрываем модальное окно
                } else {
                    setError(data.message || "Ошибка при принятии заявки");
                }
            })
            .catch((error) => {
                console.error("Error approving contract:", error);
                setError("Ошибка при принятии заявки");
            });
    };

    const handleRejectPdf = (registrationId) => {
        if (!rejectReasonPdf) {
            setError("Укажите причину отклонения");
            return;
        }
        fetch(`http://localhost:5000/admin/reject-registration/${registrationId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify({ reason: rejectReasonPdf }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === registrationId
                                ? { ...reg, status: "Отклоненный", rejectReason: rejectReasonPdf }
                                : reg
                        )
                    );
                    setRejectReasonPdf(""); // Очищаем поле причины
                    setSelectedPdfRegistrationId(null); // Закрываем модальное окно
                } else {
                    setError(data.message || "Ошибка при отклонении заявки");
                }
            })
            .catch((error) => {
                console.error("Error rejecting registration:", error);
                setError("Ошибка при отклонении заявки");
            });
    };

    if (isLoading) return <div>Загрузка...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div className="admin-approval-page">
            <h1>Проверка договора</h1>
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {registrations.map((registration) => (
                        <tr key={registration._id}>
                            <td>{registration.courseTitle}</td>
                            <td>
                                <button onClick={() => handleViewUser({
                                    username: registration.userName,
                                    email: registration.userEmail,
                                    fullname: registration.userFullname,
                                    birthdate: registration.userBirthdate,
                                    birthplace: registration.userBirthplace,
                                    education: registration.userEducation,
                                    workplace: registration.userWorkplace,
                                    jobtitle: registration.userJobtitle,
                                    homeaddress: registration.userHomeaddress,
                                    phone: registration.userPhone,
                                    passportdata: registration.userPassportdata,
                                    snils: registration.userSnils,
                                })}>
                                    {registration.userName}
                                </button>
                            </td>
                            <td>{registration.status}</td>
                            <td>
                                <button className="view-pdf-btn" onClick={() => {
                                    setSelectedPdfRegistrationId(registration._id);
                                    setPdfUrl(`http://localhost:5000/user/download-contract/${registration._id}`);
                                }}>
                                    Проверить договор
                                </button>
                                <button className="reject-btn" onClick={() => handleExpel(registration._id)}>Отчислить</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {/* Модальное окно для проверки PDF */}
            // Модальное окно для проверки PDF
            {selectedPdfRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Проверка договора</h2>
                        {/* Используем iframe для отображения PDF */}
                        <iframe
                            src={pdfUrl}
                            title="Договор"
                            width="100%"
                            height="500px"
                            frameBorder="0"
                        />
                        <div className="modal-actions">
                            <button onClick={() => handleApprovePdf(selectedPdfRegistrationId)}>Принять</button>
                            <button onClick={() => setSelectedPdfRegistrationId(null)}>Закрыть</button>
                        </div>
                    </div>
                </div>
            )}

            {/* Модальное окно для отчисления */}
            {selectedRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Укажите причину отчисления</h3>
                        <input
                            type="text"
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отчисления"
                        />
                        <button onClick={confirmExpel}>Подтвердить</button>
                        <button onClick={clearModalState}>Отмена</button>
                    </div>
                </div>
            )}
            {/* Модальное окно для просмотра информации о пользователе */}
            {selectedUser && (
                <div className="user-info-modal">
                    <div className="user-info-container">
                        <h3>Информация о пользователе</h3>
                        <table className="user-info-table">
                            <tbody>
                                <tr>
                                    <th>Логин</th>
                                    <td>{selectedUser.username || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Email</th>
                                    <td>{selectedUser.email || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>ФИО</th>
                                    <td>{selectedUser.fullname || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Дата рождения</th>
                                    <td>{selectedUser.birthdate || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место рождения</th>
                                    <td>{selectedUser.birthplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Образование</th>
                                    <td>{selectedUser.education || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место работы</th>
                                    <td>{selectedUser.workplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Должность</th>
                                    <td>{selectedUser.jobtitle || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Домашний адрес</th>
                                    <td>{selectedUser.homeaddress || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Номер телефона</th>
                                    <td>{selectedUser.phone || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Данные паспорта</th>
                                    <td>{selectedUser.passportdata || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>СНИЛС</th>
                                    <td>{selectedUser.snils || "Данные отсутствуют"}</td>
                                </tr>
                            </tbody>
                        </table>
                        <button onClick={handleCloseUserInfo}>Закрыть</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminCheckDocument;